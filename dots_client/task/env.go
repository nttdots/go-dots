package task

import "fmt"
import "github.com/nttdots/go-dots/libcoap"
import log "github.com/sirupsen/logrus"

type Env struct {
    context  *libcoap.Context
    session  *libcoap.Session
    channel  chan Event

    requests map[string] *MessageTask
}

func NewEnv(context *libcoap.Context, session *libcoap.Session) *Env {
    return &Env{
        context,
        session,
        make(chan Event, 32),
        make(map[string] *MessageTask),
    }
}

func (env *Env) Run(task Task) {
    switch t := task.(type) {
    case *MessageTask:
        key := asMapKey(t.message)
        env.requests[key] = t
    }
    go task.run(env.channel)
}

func (env *Env) HandleResponse(pdu *libcoap.Pdu) {
    key := asMapKey(pdu)
    t, ok := env.requests[key]

    if !ok {
        log.Info("Unexpected incoming PDU: %v", pdu)
    } else {
        delete(env.requests, key)
        t.stop()
        t.responseHandler(t, pdu)
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
    return fmt.Sprintf("%d[%x]", pdu.MessageID, pdu.Token)
}
