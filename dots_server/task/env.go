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
	channel       chan Event
	heartBeatList map[*libcoap.Session] *HeartBeatMessageTask
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
		make(chan Event, 32),
		make(map[*libcoap.Session] *HeartBeatMessageTask),
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
	env.channel = make(chan Event, 32)
	env.heartBeatList = make(map[*libcoap.Session] *HeartBeatMessageTask)
	return env
}

// Remove session in env
func (env *Env) RemoveSession(session *libcoap.Session) {
	delete(env.heartBeatList, session)
}

/*
 * Env running method
 * parameter:
 *  task       the task need run
 */
func (env *Env) Run(session *libcoap.Session, task Task) {
	go task.run(env.channel)
}

func (env *Env) HandleResponse(session *libcoap.Session, pdu *libcoap.Pdu) {
	t := env.heartBeatList[session]
	if !session.GetIsHeartBeatTask() {
		log.Info("Unexpected PDU: %v", pdu)
	} else if !t.isStop {
		session.SetIsHeartBeatTask(false)
		if pdu.Code == libcoap.ResponseChanged {
			session.SetIsReceiveResponseContent(true)
		}
		t.stop()
		t.responseHandler(t, pdu)
		// Reset current_missing_hb
		session.SetCurrentMissingHb(0)
	}
}

func (env *Env) HandleTimeout(session *libcoap.Session, sent *libcoap.Pdu) {
	t := env.heartBeatList[session]

	if !session.GetIsHeartBeatTask() {
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
 * Get env event
 * return:
 *  channel           the event
 */
func (env *Env) EventChannel() chan Event {
	return env.channel
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
func heartbeatTimeoutHandler(task *HeartBeatMessageTask, env *Env) {
	log.Debugf("Handle HeartBeat Timeout")
    log.Debug("Exceeded missing_hb_allowed. Stop heartbeat task...")

	// Get dot peer common name from current session
	var session *libcoap.Session
	for k, v := range env.heartBeatList {
		if v == task {
			session = k
		}
	}
    cn, err := session.DtlsGetPeerCommonName()
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
	if !session.GetIsReceiveHeartBeat() {
		// Trigger mitigation mechanism is active
		log.Debug("Start Trigger Mitigation mechanism.")
		err = controllers.TriggerMitigation(customer)
		if err != nil {
			log.WithError(err).Error("TriggerMitigation() failed")
			return
		}
	}
    session.SetIsHeartBeatTask(false)

    // log.Debugf("DTLS session: %+v has already disconnected. Terminate session...", env.session.String())
    // env.session.TerminateConnectingSession(env.context)
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
		// Get session configuration of this session by customer id
		sessionConfig, err := controllers.GetSessionConfig(customer)
		if err != nil {
			log.Errorf("[Session Mngt Thread]: Get session configuration failed.")
			return
		}
		// If DOTS server receives 2.04 but DOTS server doesn't recieve heartbeat from DOTS client,  DOTS server set 'peer-hb-status' to false
        // Else  DOTS server set 'peer-hb-status' to true
		hbValue := true
		if session.GetIsReceiveResponseContent() && !session.GetIsReceiveHeartBeat() {
			hbValue = false
		}
		hbMessage, err := messages.NewHeartBeatMessage(*session, messages.JSON_HEART_BEAT_SERVER, hbValue)
		if err != nil {
			log.Errorf("Failed to create heartbeat message")
			return
		}
		if err != nil {
			log.Errorf("[Session Mngt Thread]: Failed to create new heartbeat message. Error: %+v", err)
			return
		}
		if session.GetIsHeartBeatTask() {
			log.Debugf("[Session Mngt Thread]: Waiting for current heartbeat task (id = %+v)", env.heartBeatList[session].message.MessageID)
		} else {
			log.Debugf("[Session Mngt Thread]: Create new heartbeat message (id = %+v) to check client connection", hbMessage.MessageID)
			task := NewHeartBeatMessageTask(hbMessage, sessionConfig.MissingHbAllowedIdle,
				time.Duration(sessionConfig.AckTimeoutIdle) * time.Second,
				time.Duration(sessionConfig.HeartbeatIntervalIdle) * time.Second,
				heartbeatResponseHandler, heartbeatTimeoutHandler)
			env.heartBeatList[session] = task
			env.Run(session, task)

			session.SetIsReceiveResponseContent(false)
			session.SetIsReceiveHeartBeat(false)
		}
		time.Sleep(time.Duration(sessionConfig.HeartbeatIntervalIdle) * time.Second)
	}
}
