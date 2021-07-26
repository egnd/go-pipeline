package wpool

import "fmt"

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
