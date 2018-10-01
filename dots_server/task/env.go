package task

import (
	"time"
	"github.com/nttdots/go-dots/libcoap"
	log "github.com/sirupsen/logrus"
    dots_config "github.com/nttdots/go-dots/dots_server/config"
    "github.com/nttdots/go-dots/dots_server/controllers"
    "github.com/nttdots/go-dots/dots_server/models"
)

type Env struct {
    context         *libcoap.Context
    session         *libcoap.Session
    channel         chan Event

    pingMessageTask *PingMessageTask
}

/*
 * Env constructor
 * parameter:
 *  context       the signal context
 * return:
 *  env           the new env
 */
func NewEnv(context *libcoap.Context) *Env {
    return &Env{
        context,
        nil,
        make(chan Event, 32),
        nil,
    }
}

/*
 * The renew env with session if have
 * parameter:
 *  context        the signal context
 *  session        the current transaction session
 * return:
 *  env           the renew env
 */
func (env *Env) RenewEnv(context *libcoap.Context, currentSession *libcoap.Session) *Env {
    env.context = context
    env.session = currentSession
    env.channel = make(chan Event, 32)
    env.pingMessageTask = nil
    return env
}

/*
 * Env running method
 * parameter:
 *  task       the task need run
 */
func (env *Env) Run(task Task) {
    env.session.SetIsPingTask(true)
    switch t := task.(type) {
    case *PingMessageTask:
        env.pingMessageTask = t
    }
    go task.run(env.channel)
}

func (env *Env) HandleResponse(pdu *libcoap.Pdu) {
    t := env.pingMessageTask

    if !env.session.GetIsPingTask() {
        log.Info("Unexpected PDU: %v", pdu)
    } else {
        env.session.SetIsPingTask(false)
        t.stop()
        t.responseHandler(t, pdu)
        // Reset current_missing_hb
        env.session.SetCurrentMissingHb(0)
    }
}

func (env *Env) HandleTimeout(sent *libcoap.Pdu) {
    t := env.pingMessageTask

    if !env.session.GetIsPingTask() {
        log.Info("Unexpected PDU: %v", sent)
    } else {
        t.stop()
        
        // Couting to missing-hb
        // 0: Code of Ping task
        if sent.Code == 0 {
            env.session.SetCurrentMissingHb(env.session.GetCurrentMissingHb() + 1)
        }
        t.timeoutHandler(t, env)
    }
}

/*
 * Get env context
 * return:
 *  context           the context
 */
func (env *Env) CoapContext() *libcoap.Context {
    return env.context
}

/*
 * Get env current session
 * return:
 *  currentSession           the current session
 */
func (env *Env) CoapSession() *libcoap.Session {
    return env.session
}

/*
 * Get env event
 * return:
 *  channel           the event
 */
func (env *Env) EventChannel() chan Event {
    return env.channel
}

/*
 * Check if the missing heartbeat allowed is out
 * return:
 *  true           missing-hb-allowed is out
 *  false          missing-hb-allowed is not out
 */
func (env *Env) IsHeartbeatAllowed () bool {
    return env.session.IsHeartbeatAllowed()
}

/*
 * the response handler
 * parameter:
 *  pdu       the response model
 */
func pingResponseHandler(_ *PingMessageTask, pdu *libcoap.Pdu) {
	log.WithField("Type", pdu.Type).WithField("Code", pdu.Code).Debug("Ping Ack")
}

/*
 * the timeout handler
 * parameter:
 *  env       the env
 */
func pingTimeoutHandler(_ *PingMessageTask, env *Env) {
    log.Debugf("Ping Timeout #%+v", env.session.GetCurrentMissingHb())

	if !env.IsHeartbeatAllowed() {
		log.Debug("Exceeded missing_hb_allowed. Stop ping task...")
        
        // Get dot peer common name from current session
        cn, err := env.session.DtlsGetPeerCommonName()
        if err != nil {
            log.WithError(err).Error("DtlsGetPeercCommonName() failed")
            return
        }

        // Get customer from common name
        customer, err := models.GetCustomerByCommonName(cn)
        if err != nil || customer.Id == 0 {
            log.WithError(err).Error("Customer not found.")
            return
        }

        // Trigger mitigation mechanism is active
        log.Debug("Start Trigger Mitigation mechanism.")
        err = controllers.TriggerMitigation(customer)
        if err != nil {
            log.WithError(err).Error("TriggerMitigation() failed")
            return
        }
        env.session.SetIsPingTask(false)

        log.Debugf("DTLS session: %+v has already disconnected. Terminate session...", env.session.String())
        env.session.TerminateConnectingSession(env.context)
        return
    }
    log.Debugf("[Session Mngt Thread]: Re-send ping message (id = %+v) to check client connection", env.pingMessageTask.message.MessageID)
    env.Run(env.pingMessageTask)
}

/*
 * the monitoring thread that manage session traffic from dots client every cycle
 * parameter:
 *  pdu       the response model
 */
func (env *Env) ManageSessionTraffic() {
    config := dots_config.GetServerSystemConfig()

    for { 
        // loop on sessions in map that connecting to server
        time.Sleep(time.Duration(config.SessionInterval) * time.Second)
        log.Debugf("[Session Mngt Thread]: The number of connecting session: %+v", len(libcoap.ConnectingSessions()))
        for _, session := range libcoap.ConnectingSessions() {
            if session.GetIsAlive() == false {
                if session.GetIsPingTask() == true {
                    log.Debugf("[Session Mngt Thread]: Waiting for current ping task (id = %+v)", env.pingMessageTask.message.MessageID)
                    break
                }

                // Get dot peer common name from current session
                cn, err := session.DtlsGetPeerCommonName()
                if err != nil {
                    log.WithError(err).Error("[Session Mngt Thread]: DtlsGetPeercCommonName() failed")
                    return
                }

                // Get customer from common name
                customer, err := models.GetCustomerByCommonName(cn)
                if err != nil || customer.Id == 0 {
                    log.WithError(err).Error("[Session Mngt Thread]: Customer not found.")
                    return
                }

                // Get session configuration of this session by customer id
                sessionConfig, err := controllers.GetSessionConfig(customer)
                if err != nil {
                    log.Errorf("[Session Mngt Thread]: Get session configuration failed.")
                }

                // Update session configuration and start to run ping task
                session.SetSessionConfig(sessionConfig.MissingHbAllowed, sessionConfig.MaxRetransmit, sessionConfig.AckTimeout, sessionConfig.AckRandomFactor)
                env.session = session

                // new message ping to sent to client
                pingMessage := newPingMessage(env)
                log.Debugf("[Session Mngt Thread]: Create new ping message (id = %+v) to check client connection", pingMessage.MessageID)
                env.Run(NewPingMessageTask(pingMessage, session.GetCurrentMissingHb(), time.Duration(0), time.Duration(0),
                    pingResponseHandler, pingTimeoutHandler ))
            } else {
                // Set false for checking session at next cycle 
                log.Debug("[Session Mngt Thread]: Skip ping task message. Dots client has just negotiated recently.")
                session.SetIsAlive(false)
            }
        }

    }
}
