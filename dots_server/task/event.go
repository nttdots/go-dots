package task

type Event interface {
    Task() Task
    Handle(*Env)
}

type EventBase struct {
    task Task
}

func (e EventBase) Task() Task {
    return e.task
}
