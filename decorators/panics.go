package decorators

import (
	"fmt"

	"github.com/egnd/go-pipeline"
)

// CatchPanic is catching panic and return it as an error.
func CatchPanic(next pipeline.TaskExecutor) pipeline.TaskExecutor {
	return func(task pipeline.Task) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
			}
		}()

		err = next(task)

		return
	}
}

// ThrowPanic is throws task error as a panic.
func ThrowPanic(next pipeline.TaskExecutor) pipeline.TaskExecutor {
	return func(task pipeline.Task) error {
		if err := next(task); err != nil {
			panic(err)
		}

		return nil
	}
}
