package task

import "time"
import "github.com/nttdots/go-dots/libcoap"

type ResponseHandler func(*MessageTask, *libcoap.Pdu)
type TimeoutHandler  func(*MessageTask)

type MessageTask struct {
    TaskBase

    message  *libcoap.Pdu

    interval time.Duration
    retry    int
    timeout  time.Duration

    responseHandler ResponseHandler
    timeoutHandler  TimeoutHandler
}

type TimeoutEvent struct { EventBase }
type MessageEvent struct { EventBase }

func NewMessageTask(message *libcoap.Pdu,
                    interval time.Duration,
                    retry int,
                    timeout time.Duration,
                    responseHandler ResponseHandler,
                    timeoutHandler TimeoutHandler) *MessageTask {
    return &MessageTask {
        newTaskBase(),
        message,
        interval,
        retry,
        timeout,
        responseHandler,
        timeoutHandler,
    }
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

    select {
    case <- t.stopChan:
        return
    // case <- timeout:
    //     out <- &TimeoutEvent{ EventBase{ t } }
    //     t.stop()
    }
}

func (e *MessageEvent) Handle(env *Env) {
    task := e.Task().(*MessageTask)
    env.session.Send(task.message)
}

func (e *TimeoutEvent) Handle(env *Env) {
    task := e.Task().(*MessageTask)
    task.timeoutHandler(task)
}
