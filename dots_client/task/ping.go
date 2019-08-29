package task

import "time"
import "github.com/nttdots/go-dots/libcoap"
import log "github.com/sirupsen/logrus"

type PingResponseHandler func(*PingTask, *libcoap.Pdu)
type PingTimeoutHandler  func(*PingTask, *Env)

type PingTask struct {
    TaskBase

    interval        time.Duration
    responseHandler PingResponseHandler
    timeoutHandler  PingTimeoutHandler
    current_ping_id  string
}

type PingEvent struct { EventBase }



func NewPingTask(interval time.Duration, responseHandler PingResponseHandler, timeoutHandler PingTimeoutHandler) *PingTask {
    return &PingTask {
        newTaskBase(),
        interval,
        responseHandler,
        timeoutHandler,
        "",
    }
}

func (t *PingTask) run(out chan Event) {
    for {
        select {
        case <- t.stopChan:{
            log.Debug("Current ping task ended.")
            return
        }
        case <- time.After(t.interval):
            out <- &PingEvent{ EventBase{ t } }
        }
    }
}

func newPingMessage(env *Env) *libcoap.Pdu {
    pdu := &libcoap.Pdu{}
    pdu.Type = libcoap.TypeCon
    pdu.Code = 0
    pdu.MessageID = env.session.NewMessageID()
    return pdu
}

func (e *PingEvent) Handle(env *Env) {
    pingTask := e.task.(*PingTask)
    currentTask := env.requests[pingTask.current_ping_id]

    if currentTask != nil {
        log.Debugf("Waiting for current ping message to be completed...")
        return
    }

    // Send new ping message
    pdu := newPingMessage(env)
    task := e.Task().(*PingTask)

    newTask := NewMessageTask(
        pdu,
        time.Duration(0),
        0,
        time.Duration(0),
        false,
        func (_ *MessageTask, pdu *libcoap.Pdu, env *Env) {
            task.responseHandler(task, pdu)
        },
        func (*MessageTask, *Env) {
            task.timeoutHandler(task, env)
        })
    env.Run(newTask)
    pingTask.current_ping_id = newTask.message.AsMapKey()
    log.Debugf ("Sent new ping message (id = %+v)", newTask.message.MessageID )
}

func (t * PingTask) IsRunnable() bool {
    return t.interval > 0
}