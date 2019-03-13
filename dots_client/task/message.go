package task

import "time"
import "github.com/nttdots/go-dots/libcoap"
import log "github.com/sirupsen/logrus"

type ResponseHandler func(*MessageTask, *libcoap.Pdu)
type TimeoutHandler  func(*MessageTask, map[string] *MessageTask)

type MessageTask struct {
    TaskBase

    message  *libcoap.Pdu
    response chan *libcoap.Pdu

    interval time.Duration
    retry    int
    timeout  time.Duration

    isStop bool
    responseHandler ResponseHandler
    timeoutHandler  TimeoutHandler
}

type TimeoutEvent struct { EventBase }
type MessageEvent struct { EventBase }

func NewMessageTask(message *libcoap.Pdu,
                    interval time.Duration,
                    retry int,
                    timeout time.Duration,
                    isStop bool,
                    responseHandler ResponseHandler,
                    timeoutHandler TimeoutHandler) *MessageTask {
    return &MessageTask {
        newTaskBase(),
        message,
        make(chan *libcoap.Pdu),
        interval,
        retry,
        timeout,
        isStop,
        responseHandler,
        timeoutHandler,
    }
}

func (task *MessageTask) GetMessage() (*libcoap.Pdu) {
    return task.message
}

func (task *MessageTask) SetMessage(pdu *libcoap.Pdu) {
    task.message = pdu
}

func (t *MessageTask) run(out chan Event) {
    timeout := time.After(t.timeout)

    out <- &MessageEvent{ EventBase{ t } }

    for i := 0; i < t.retry; i++ {
        select {
        case <- t.stopChan:
            return
        case <- time.After(t.interval):
            out <- &MessageEvent{ EventBase{ t } }
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
            log.Debug("Mitigation request timeout")
            t.isStop = true
            out <- &TimeoutEvent{ EventBase{ t } }
            t.stop()
        }
    }else {
        select {
        case <- t.stopChan:
            return
        }
    }
}

func (e *MessageEvent) Handle(env *Env) {
    task := e.Task().(*MessageTask)
    env.session.Send(task.message)
}

func (e *TimeoutEvent) Handle(env *Env) {
    task := e.Task().(*MessageTask)
    task.timeoutHandler(task, env.requests)
}

func (t *MessageTask) AddResponse(pdu *libcoap.Pdu) {
    t.response <- pdu
}