package wpool_test

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/egnd/go-wpool/v2"
	"github.com/egnd/go-wpool/v2/interfaces"
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
		t.Run(fmt.Sprint(k), func(tt *testing.T) {
			pipeline := make(chan interfaces.Worker)
			pool := wpool.NewPipelinePool(pipeline, nil)

			for i := 0; i <= test.workersCnt; i++ {
				pool.AddWorker(wpool.NewPipelineWorker(pipeline))
			}

			var wg sync.WaitGroup
			for i := 0; i <= test.tasksCnt; i++ {
				wg.Add(1)
				task := &interfaces.MockTask{}
				task.On("Do").After(time.Duration(rand.Intn(10)) * time.Millisecond).Run(func(_ mock.Arguments) { wg.Done() }).Once()
				defer task.AssertExpectations(tt)
				assert.NoError(tt, pool.AddTask(task))
			}

			wg.Wait()
			assert.NoError(tt, pool.Close())
		})
	}
}

func Test_PipelinePool_Workers_Errors(t *testing.T) {
	pipeline := make(chan interfaces.Worker)

	logger := interfaces.MockLogger{}
	defer logger.AssertExpectations(t)
	logger.On("Errorf", mock.Anything, mock.Anything, "taskid")

	worker := interfaces.MockWorker{}
	defer worker.AssertExpectations(t)
	worker.On("Do", mock.Anything).Return(errors.New("error"))
	worker.On("Close", mock.Anything).Return(errors.New("error"))

	task := interfaces.MockTask{}
	defer task.AssertExpectations(t)
	task.On("GetID").Return("taskid").Maybe()

	pool := wpool.NewPipelinePool(pipeline, &logger)

	go func() {
		pipeline <- &worker
	}()

	pool.AddWorker(&worker)
	assert.NoError(t, pool.AddTask(&task))
	assert.EqualError(t, pool.Close(), "error")
}

func Test_PipelinePool_Close_Error(t *testing.T) {
	pipeline := make(chan interfaces.Worker)
	pool := wpool.NewPipelinePool(pipeline, nil)
	assert.NoError(t, pool.Close())
	assert.EqualError(t, pool.AddTask(nil), "pool is closed")
}
