package task

import "time"
import "github.com/nttdots/go-dots/libcoap"
import log "github.com/sirupsen/logrus"

type PingResponseHandler func(*PingMessageTask, *libcoap.Pdu)
type PingTimeoutHandler  func(*PingMessageTask, *Env)

type PingMessageTask struct {
    TaskBase

    message  *libcoap.Pdu

    retry    int
    timeout  time.Duration
    interval time.Duration

    responseHandler PingResponseHandler
    timeoutHandler  PingTimeoutHandler
}

type TimeoutEvent struct { EventBase }
type PingEvent struct { EventBase }

func NewPingMessageTask(message *libcoap.Pdu, retry int, timeout time.Duration, interval time.Duration,
     responseHandler PingResponseHandler, timeoutHandler PingTimeoutHandler) *PingMessageTask {
    return &PingMessageTask {
        newTaskBase(),
        message,
        retry,
        timeout,
        interval,
        responseHandler,
        timeoutHandler,
    }
}

func (t *PingMessageTask) run(out chan Event) {
    timeout := time.After(t.timeout)
    out <- &PingEvent{ EventBase{ t } }

    for i := 0; i < t.retry; i++ {
        select {
        case <- t.stopChan:
            log.Debug("Current ping task ended.")
            return
        case <- time.After(t.interval):
            out <- &PingEvent{ EventBase{ t } }
        case <- timeout:
            out <- &TimeoutEvent{ EventBase{ t } }
            t.stop()
        }
    }

    select {
    case <- t.stopChan:
        return
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
    task := e.Task().(*PingMessageTask)
    env.session.Send(task.message)
}

func (e *TimeoutEvent) Handle(env *Env) {
    task := e.Task().(*PingMessageTask)
    task.timeoutHandler(task, env)
}