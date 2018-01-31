package task

import "time"
import "github.com/nttdots/go-dots/libcoap"

type PingResponseHandler func(*PingTask, *libcoap.Pdu)
type PingTimeoutHandler  func(*PingTask)

type PingTask struct {
    TaskBase

    interval        time.Duration
    responseHandler PingResponseHandler
    timeoutHandler  PingTimeoutHandler
}

type PingEvent struct { EventBase }

func NewPingTask(interval time.Duration, responseHandler PingResponseHandler, timeoutHandler PingTimeoutHandler) *PingTask {
    return &PingTask {
        newTaskBase(),
        interval,
        responseHandler,
        timeoutHandler,
    }
}

func (t *PingTask) run(out chan Event) {
    for {
        select {
        case <- t.stopChan:
            return
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
    pdu := newPingMessage(env)
    task := e.Task().(*PingTask)

    newTask := NewMessageTask(
        pdu,
        time.Duration(0),
        0,
        time.Duration(task.interval * 2),
        func (_ *MessageTask, pdu *libcoap.Pdu) {
            task.responseHandler(task, pdu)
        },
        func (*MessageTask) {
            task.timeoutHandler(task)
        })
    env.Run(newTask)
}
