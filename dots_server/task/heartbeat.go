package task

import (
    "time"
    "github.com/nttdots/go-dots/libcoap"
)

type HeartBeatResponseHandler func(*HeartBeatMessageTask, *libcoap.Pdu)
type HeartBeatTimeoutHandler  func(*HeartBeatMessageTask, *Env)

type HeartBeatMessageTask struct {
    TaskBase

    message  *libcoap.Pdu

    retry    int
    timeout  time.Duration
    interval time.Duration

    responseHandler HeartBeatResponseHandler
    timeoutHandler  HeartBeatTimeoutHandler
}

type TimeoutEvent struct { EventBase }
type HeartBeatEvent struct { EventBase }

func NewHeartBeatMessageTask(message *libcoap.Pdu, retry int, timeout time.Duration, interval time.Duration,
     responseHandler HeartBeatResponseHandler, timeoutHandler HeartBeatTimeoutHandler) *HeartBeatMessageTask {
    return &HeartBeatMessageTask {
        newTaskBase(),
        message,
        retry,
        timeout,
        interval,
        responseHandler,
        timeoutHandler,
    }
}

func (t *HeartBeatMessageTask) run(out chan Event) {
    timeout := time.After(t.timeout)
    out <- &HeartBeatEvent{ EventBase{ t } }

    for i := 0; i < t.retry; i++ {
        select {
        case <- t.stopChan:{
            return
        }
        case <- time.After(t.interval):
            out <- &HeartBeatEvent{ EventBase{ t } }
        case <- timeout:
            out <- &TimeoutEvent{ EventBase{ t } }
            t.stop()
        }
    }
    if t.message.Type == libcoap.TypeNon {
        select {
        case <- t.stopChan:
            return
        case <- timeout:
            out <- &TimeoutEvent{ EventBase{ t } }
            t.stop()
        }
    } else {
        select {
        case <- t.stopChan:
            return
        }
    }
}

func (e *HeartBeatEvent) Handle(env *Env) {
    task := e.Task().(*HeartBeatMessageTask)
    env.session.SetIsHeartBeatTask(true)
    env.session.Send(task.message)
}

func (e *TimeoutEvent) Handle(env *Env) {
    task := e.Task().(*HeartBeatMessageTask)
    task.timeoutHandler(task, env)
}