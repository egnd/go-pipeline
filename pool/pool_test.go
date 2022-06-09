package pool_test

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/egnd/go-pipeline"
	"github.com/egnd/go-pipeline/mocks"
	"github.com/egnd/go-pipeline/pool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_PipelinePool(t *testing.T) {
	cases := []struct {
		tasksCnt   int
		workersCnt int
	}{
		{
			workersCnt: 1,
			tasksCnt:   10,
		},
		{
			workersCnt: 2,
			tasksCnt:   21,
		},
		{
			workersCnt: 10,
			tasksCnt:   502,
		},
	}
	for k, test := range cases {
		t.Run(fmt.Sprint(k+1), func(tt *testing.T) {
			test := test

			bus := make(chan pipeline.Doer)

			workers := []pipeline.Doer{}
			for i := 0; i <= test.workersCnt; i++ {
				workers = append(workers, pool.NewWorker(bus))
			}

			pipe := pool.NewPool(bus, workers...)

			var wg sync.WaitGroup
			for i := 0; i <= test.tasksCnt; i++ {
				task := &mocks.Task{}
				task.On("Do").After(time.Duration(rand.Intn(10)) * time.Millisecond).Run(func(_ mock.Arguments) { wg.Done() }).Once().Return(nil)
				defer task.AssertExpectations(tt)

				wg.Add(1)
				assert.NoError(tt, pipe.Push(task))
			}

			wg.Wait()
			assert.NoError(tt, pipe.Close())
		})
	}
}

func Test_Pool_CLose_Errors(t *testing.T) {
	bus := make(chan pipeline.Doer)
	worker := mocks.Doer{}
	defer worker.AssertExpectations(t)
	worker.On("Close", mock.Anything).Return(errors.New("error")).Once()
	pipe := pool.NewPool(bus, &worker)
	assert.EqualError(t, pipe.Close(), "error")
	assert.EqualError(t, pipe.Close(), "pool close err: close of closed channel")
}

// func Test_Pool_Push_Errors(t *testing.T) { @TODO:
// 	worker := mocks.Doer{}
// 	worker.On("Do", mock.Anything).Return(errors.New("worker error"))
// 	defer worker.AssertExpectations(t)

// 	bus := make(chan pipeline.Doer, 1)
// 	bus <- &worker

// 	assert.Panics(t, func() {
// 		pool.NewPool(bus, &worker).Push(nil)
// 	})
// }

func Test_PipelinePool_Close_Error(t *testing.T) {
	bus := make(chan pipeline.Doer)
	pipe := pool.NewPool(bus)
	assert.NoError(t, pipe.Close())
	assert.EqualError(t, pipe.Push(nil), "pool push err: send on closed channel")
}
