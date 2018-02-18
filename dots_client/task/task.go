package task

type Task interface {
    stop()
    run(chan Event)
}

type TaskBase struct {
    stopChan chan int
}

func newTaskBase() TaskBase {
    return TaskBase{ make(chan int) }
}

func (t TaskBase) stop() {
    t.stopChan <- 0
}
