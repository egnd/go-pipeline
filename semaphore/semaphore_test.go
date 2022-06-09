package semaphore_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/egnd/go-pipeline/mocks"
	"github.com/egnd/go-pipeline/semaphore"
	"github.com/stretchr/testify/assert"
)

func Test_Semaphore(t *testing.T) {
	cases := []struct {
		buffSize int
		tasks    []mocks.Task
	}{
		{
			buffSize: 2,
			tasks: func() (res []mocks.Task) {
				for i := 0; i < 20; i++ {
					task := mocks.Task{}
					task.On("Do").Once().After(time.Duration(rand.Intn(10)) * time.Millisecond).Return(nil)
					res = append(res, task)
				}
				return
			}(),
		},
	}

	for k, test := range cases {
		t.Run(fmt.Sprint(k+1), func(t *testing.T) {
			pipe := semaphore.NewSemaphore(test.buffSize)

			for _, task := range test.tasks {
				task := task
				pipe.Push(&task)
				defer task.AssertExpectations(t)
			}

			pipe.Close()
		})
	}
}

func Test_Semaphore_Errors(t *testing.T) {
	pipe := semaphore.NewSemaphore(3)
	assert.NoError(t, pipe.Close())
	assert.EqualError(t, pipe.Close(), "semaphore close err: close of closed channel")
	assert.EqualError(t, pipe.Push(&mocks.Task{}), "semaphore do err: send on closed channel")
}
