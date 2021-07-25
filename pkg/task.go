package pkg

import "sync"

// ITask is a task interface.
type ITask interface {
	Name() string
	Do() error
}

// ITask is a task struct.
type Task struct {
	Title    string
	Callback func() error
	Wg       *sync.WaitGroup
}

// Name is returning task name.
func (t *Task) Name() string {
	return t.Title
}

// Do is executing task logic.
func (t *Task) Do() error {
	if t.Wg != nil {
		defer t.Wg.Done()
	}

	if t.Callback != nil {
		return t.Callback()
	}

	return nil
}
