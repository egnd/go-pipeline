package wpool

import "sync"

// ITask is a task interface.
type ITask interface {
	GetName() string
	Do() error
}

// Task is a task struct.
type Task struct {
	Name     string
	Callback func(*Task) error
	Wg       *sync.WaitGroup
}

// GetName is returning task name.
func (t *Task) GetName() string {
	return t.Name
}

// Do is executing task logic.
func (t *Task) Do() error {
	if t.Wg != nil {
		defer t.Wg.Done()
	}

	if t.Callback != nil {
		return t.Callback(t)
	}

	return nil
}
