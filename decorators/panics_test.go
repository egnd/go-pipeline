package decorators_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/egnd/go-pipeline"
	"github.com/egnd/go-pipeline/decorators"
	"github.com/egnd/go-pipeline/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_CatchPanic(t *testing.T) {
	cases := []struct {
		panic string
	}{
		{"123"},
		{"dv78v9s"},
		{},
	}

	for k, test := range cases {
		t.Run(fmt.Sprint(k+1), func(t *testing.T) {
			task := &mocks.Task{}
			if test.panic == "" {
				task.On("Do").Return(nil)
				pipeline.NewTaskExecutor([]pipeline.TaskDecorator{func(next pipeline.TaskExecutor) pipeline.TaskExecutor {
					return func(task pipeline.Task) (err error) {
						assert.NoError(t, next(task))
						return
					}
				}, decorators.CatchPanic})(task)
			} else {
				task.On("Do").Panic(test.panic)
				pipeline.NewTaskExecutor([]pipeline.TaskDecorator{func(next pipeline.TaskExecutor) pipeline.TaskExecutor {
					return func(task pipeline.Task) (err error) {
						assert.EqualValues(t, errors.New(test.panic), next(task))
						return
					}
				}, decorators.CatchPanic})(task)
			}
		})
	}
}

func Test_ThrowPanic(t *testing.T) {
	cases := []struct {
		err error
	}{
		{},
		{errors.New("test error")},
	}

	for k, test := range cases {
		t.Run(fmt.Sprint(k+1), func(t *testing.T) {
			task := &mocks.Task{}
			task.On("Do").Return(test.err)

			if test.err == nil {
				pipeline.NewTaskExecutor([]pipeline.TaskDecorator{decorators.ThrowPanic})(task)
			} else {
				assert.PanicsWithError(t, test.err.Error(), func() {
					pipeline.NewTaskExecutor([]pipeline.TaskDecorator{decorators.ThrowPanic})(task)
				})
			}
		})
	}
}
