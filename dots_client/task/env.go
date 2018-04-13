package task

import (
	"github.com/shopspring/decimal"
	"fmt"
	"github.com/nttdots/go-dots/libcoap"
	log "github.com/sirupsen/logrus"
	"reflect"
)

type Env struct {
    context  *libcoap.Context
    session  *libcoap.Session
    channel  chan Event

    requests map[string] *MessageTask
    missing_hb_allowed int
    current_missing_hb int
    pingTask *PingTask
}

func NewEnv(context *libcoap.Context, session *libcoap.Session) *Env {
    return &Env{
        context,
        session,
        make(chan Event, 32),
        make(map[string] *MessageTask),
        0,
        0,
        nil,
    }
}

func (env *Env) RenewEnv(context *libcoap.Context, session *libcoap.Session) *Env {
    env.context = context
    env.session = session
    env.channel = make(chan Event, 32)
    env.requests = make(map[string] *MessageTask)
    env.current_missing_hb = 0
    env.pingTask = nil
    return env
}

func (env *Env) SetRetransmitParams(maxRetransmit int, ackTimeout int, ackRandomFactor decimal.Decimal){
    env.session.SetMaxRetransmit(maxRetransmit)
    env.session.SetAckTimeout(ackTimeout)
    env.session.SetAckRandomFactor(ackRandomFactor)
}

func (env *Env) SetMissingHbAllowed(missing_hb_allowed int) {
    env.missing_hb_allowed = missing_hb_allowed
}


func (env *Env) Run(task Task) {
    if (reflect.TypeOf(task) == reflect.TypeOf(&PingTask{})) && (!task.(*PingTask).IsRunnable()) {
        log.Debug("Ping task is disabled. Do not start ping task.")
        return
    }

    switch t := task.(type) {
    case *MessageTask:
        key := asMapKey(t.message)
        env.requests[key] = t

    case *PingTask:
        env.pingTask = t
    }
    go task.run(env.channel)
}

func (env *Env) HandleResponse(pdu *libcoap.Pdu) {
    key := asMapKey(pdu)
    t, ok := env.requests[key]

    if !ok {
        log.Info("Unexpected incoming PDU: %v", pdu)
    } else {
        log.Debugf("Success incoming PDU(HandleResponse): %+v", pdu)
        delete(env.requests, key)
        t.stop()
        t.responseHandler(t, pdu)
        // Reset current_missing_hb
	    env.current_missing_hb = 0
    }
}

func (env *Env) HandleTimeout(sent *libcoap.Pdu) {
    key := asMapKey(sent)
    t, ok := env.requests[key]

    if !ok {
        log.Info("Unexpected PDU: %v", sent)
    } else {
        delete(env.requests, key)
        t.stop()

        // Couting to missing-hb
        env.current_missing_hb = env.current_missing_hb + 1

        t.timeoutHandler(t)
    }
}

func (env *Env) CoapContext() *libcoap.Context {
    return env.context
}

func (env *Env) CoapSession() *libcoap.Session {
    return env.session
}

func (env *Env) EventChannel() chan Event {
    return env.channel
}

func asMapKey(pdu *libcoap.Pdu) string {
    // return fmt.Sprintf("%d[%x]", pdu.MessageID, pdu.Token)
    // return fmt.Sprintf("%x", pdu.Token)
    return fmt.Sprintf("%d", pdu.MessageID)
}

func (env *Env) IsHeartbeatAllowed () bool {
    return env.current_missing_hb < env.missing_hb_allowed
}

func (env *Env) StopPing() {
    env.pingTask.stop()
}

func (env *Env) CurrentMissingHb() int {
    return env.current_missing_hb
}
