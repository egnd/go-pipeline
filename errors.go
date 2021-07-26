package wpool

import (
	"fmt"
)

// ErrIsClosed is an error which shows that something is closed for writing.
type ErrIsClosed struct {
	EntityName string
}

func (e ErrIsClosed) Error() string {
	return fmt.Sprintf("%s is closed", e.EntityName)
}

// ErrPanic is an error for recovered panics handling.
type ErrPanic struct {
	Data interface{}
}

func (e ErrPanic) Error() string {
	return fmt.Sprintf("panic occurred: %v", e.Data)
}

// ErrWrapper is a wrapper for errors.
type ErrWrapper struct {
	Msg string
	Err error
}

func (e ErrWrapper) Error() string {
	return fmt.Sprintf("%s: %s", e.Msg, e.Err.Error())
}

// ErrTaskTimeout is an error for task execution timeout.
type ErrTaskTimeout struct {
	TaskName string
}

// Timeout shows if is the error a timeout.
func (e *ErrTaskTimeout) Timeout() bool {
	return true
}

// Temporary shows if is the error temporary.
func (e *ErrTaskTimeout) Temporary() bool {
	return true
}

func (e *ErrTaskTimeout) Error() string {
	return fmt.Sprintf("task timeout: %s", e.TaskName)
}
