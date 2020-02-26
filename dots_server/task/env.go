package task

import (
	"time"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/controllers"
	"github.com/nttdots/go-dots/dots_server/models"
	"github.com/nttdots/go-dots/libcoap"
	log "github.com/sirupsen/logrus"
)

type Env struct {
	context       *libcoap.Context
	session       *libcoap.Session
	channel       chan Event
	hbMessageTask *HeartBeatMessageTask
}

var env *Env

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
	env.hbMessageTask = nil
	return env
}

/*
 * Env running method
 * parameter:
 *  task       the task need run
 */
func (env *Env) Run(task Task) {
	switch t := task.(type) {
	case *HeartBeatMessageTask:
		env.hbMessageTask = t
	}
	go task.run(env.channel)
}

func (env *Env) HandleResponse(pdu *libcoap.Pdu) {
	t := env.hbMessageTask

	if !env.session.GetIsHeartBeatTask() {
		log.Info("Unexpected PDU: %v", pdu)
	} else {
		env.session.SetIsHeartBeatTask(false)
		if pdu.Code == libcoap.ResponseChanged {
			env.session.SetIsReceiveResponseContent(true)
		}
		t.stop()
		t.responseHandler(t, pdu)
		// Reset current_missing_hb
		env.session.SetCurrentMissingHb(0)
	}
}

func (env *Env) HandleTimeout(sent *libcoap.Pdu) {
	t := env.hbMessageTask

	if !env.session.GetIsHeartBeatTask() {
		log.Info("Unexpected PDU: %v", sent)
	} else {
		t.stop()
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

// Set env session
func (env *Env) SetCoapSession(session *libcoap.Session) {
	env.session = session
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
func (env *Env) IsHeartbeatAllowed() bool {
	return env.session.IsHeartbeatAllowed()
}

/*
 * Get env
 * return:
 *  env           the evironment
 */
func SetEnv(value *Env) {
	env = value
}

/*
 * Get env
 * return:
 *  env           the evironment
 */
func GetEnv() *Env {
	return env
}

/*
 * Get env heartbeat message task
 * return:
 *  hbMessageTask           the heartbeat mesage task
 */
func (env *Env) GetHbMessageTask() *HeartBeatMessageTask {
	return env.hbMessageTask
}

/*
 * the response handler
 * parameter:
 *  pdu       the response model
 */
func heartbeatResponseHandler(t *HeartBeatMessageTask, pdu *libcoap.Pdu) {
	log.WithField("Type", pdu.Type).WithField("Code", pdu.Code).Debug("HeartBeat")
	if pdu.Code != libcoap.ResponseChanged {
		log.Debugf("Error message: %+v", string(pdu.Data))
	}
}

/*
 * the timeout handler
 * parameter:
 *  env       the env
 */
func heartbeatTimeoutHandler(_ *HeartBeatMessageTask, env *Env) {
	log.Debugf("HeartBeat Timeout")
    log.Debug("Exceeded missing_hb_allowed. Stop heartbeat task...")

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

	// If DOTS server heartbeat timeout and doesn't receive heartbeat from DOTS client, DOTS server will start trigger mitigation
	if env.CoapSession().GetIsReceiveHeartBeat() == false {
		// Trigger mitigation mechanism is active
		log.Debug("Start Trigger Mitigation mechanism.")
		err = controllers.TriggerMitigation(customer)
		if err != nil {
			log.WithError(err).Error("TriggerMitigation() failed")
			return
		}
	}
    env.session.SetIsHeartBeatTask(false)

    // log.Debugf("DTLS session: %+v has already disconnected. Terminate session...", env.session.String())
    // env.session.TerminateConnectingSession(env.context)
    return
}

// Handle heartbeat from server to client
func (env *Env) HeartBeatMechaism(session *libcoap.Session, customer *models.Customer) {
	// Set isSentHeartBeat is true to check the DOTS server sent ping to DOTS client
	session.SetIsSentHeartBeat(true)
	for {
		// If session is closed, DOTS server will doesn't sent Ping to DOTS client
		sessions := libcoap.ConnectingSessions()
		if sessions[session.GetSessionPtr()] == nil {
			return
		}
		env.session = session
		// Get session configuration of this session by customer id
		sessionConfig, err := controllers.GetSessionConfig(customer)
		if err != nil {
			log.Errorf("[Session Mngt Thread]: Get session configuration failed.")
			return
		}
		// If DOTS server receives 2.04 but DOTS server doesn't recieve heartbeat from DOTS client,  DOTS server set 'peer-hb-status' to false
        // Else  DOTS server set 'peer-hb-status' to true
		hbValue := true
		if session.GetIsReceiveResponseContent() == true && session.GetIsReceiveHeartBeat() == false {
			hbValue = false
		}
		hbMessage, err := messages.NewHeartBeatMessage(*env.CoapSession(), messages.JSON_HEART_BEAT_SERVER, hbValue)
		if err != nil {
			log.Errorf("Failed to create heartbeat message")
			return
		}
		if err != nil {
			log.Errorf("[Session Mngt Thread]: Failed to create new heartbeat message. Error: %+v", err)
			return
		}
		if env.CoapSession().GetIsHeartBeatTask() == true {
			log.Debugf("[Session Mngt Thread]: Waiting for current heartbeat task (id = %+v)", env.hbMessageTask.message.MessageID)
		} else {
			log.Debugf("[Session Mngt Thread]: Create new heartbeat message (id = %+v) to check client connection", hbMessage.MessageID)
			env.Run(NewHeartBeatMessageTask(hbMessage, sessionConfig.MissingHbAllowedIdle,
				time.Duration(sessionConfig.AckTimeoutIdle) * time.Second,
				time.Duration(sessionConfig.HeartbeatIntervalIdle) * time.Second,
				heartbeatResponseHandler, heartbeatTimeoutHandler))

			session.SetIsReceiveResponseContent(false)
			session.SetIsReceiveHeartBeat(false)
		}
		time.Sleep(time.Duration(sessionConfig.HeartbeatIntervalIdle) * time.Second)
	}
}
