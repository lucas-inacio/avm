package manager

import (
	"sync"
)

type Progress interface {
	sync.Locker
	Done() chan bool
	GetProgress() float32

	SetDone()
	SetProgress(value int)

	GetError() error
	SetError(err error)

	GetTotal() int
}

type TaskProgress struct {
	sync.Mutex
	progress int
	done chan bool
	err error
	total int
}

func NewTaskProgress(total int) *TaskProgress {
	return &TaskProgress{
		progress: 0.0,
		done: make(chan bool),
		err: nil,
		total: total,
	}
}

func (task *TaskProgress) GetProgress() int {
	task.Lock()
	defer task.Unlock()

	return task.progress
}

func (task *TaskProgress) SetProgress(value int) {
	task.Lock()
	defer task.Unlock()

	task.progress = value
}

func (task *TaskProgress) Done() chan bool {
	return task.done
}

func (task *TaskProgress) SetDone() {
	task.done <- true
}

func (task *TaskProgress) GetError() error {
	return task.err
}

func (task *TaskProgress) SetError(err error) {
	task.Lock()
	defer task.Unlock()
	task.err = err
}

func (task *TaskProgress) GetTotal() int {
	return task.total
}