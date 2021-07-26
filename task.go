package wpool

import (
	"context"
	"sync"
)

// ITask is a task interface.
type ITask interface {
	GetName() string
	Do(context.Context) error
}

// Task is a task struct.
type Task struct {
	Name     string
	Callback func(context.Context, *Task) error
	Wg       *sync.WaitGroup
}

// GetName is returning task name.
func (t *Task) GetName() string {
	return t.Name
}

// Do is executing task logic.
func (t *Task) Do(ctx context.Context) error {
	if t.Wg != nil {
		defer t.Wg.Done()
	}

	if t.Callback != nil {
		return t.Callback(ctx, t)
	}

	return nil
}
