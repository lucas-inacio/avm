package manager

import (
	"sync"
)

type Progress interface {
	sync.Locker
	Done() chan bool
	GetProgress() float32

	SetDone()
	SetProgress(value float32)

	GetError() error
	SetError(err error)
}

type TaskProgress struct {
	sync.Mutex
	progress float32
	done chan bool
	err error
}

func NewTaskProgress() *TaskProgress {
	return &TaskProgress{
		progress: 0.0,
		done: make(chan bool),
		err: nil,
	}
}

func (task *TaskProgress) GetProgress() float32 {
	task.Lock()
	defer task.Unlock()

	return task.progress
}

func (task *TaskProgress) SetProgress(value float32) {
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