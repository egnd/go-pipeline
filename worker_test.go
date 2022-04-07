package wpool_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/egnd/go-wpool"
	"github.com/egnd/go-wpool/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Worker(t *testing.T) {
	cases := []struct {
		cfg      wpool.WorkerCfg
		tasksCnt int
		err      error
	}{
		{
			cfg: wpool.WorkerCfg{
				TasksChanBuff: 10,
				TaskTTL:       time.Minute,
				Pipeline:      make(chan<- wpool.IWorker, 100),
			},
			tasksCnt: 21,
		},
		{
			cfg:      wpool.WorkerCfg{TasksChanBuff: 20},
			tasksCnt: 10,
			err:      errors.New("error"),
		},
	}
	ctx := context.Background()
	for k, test := range cases {
		t.Run(fmt.Sprint(k), func(tt *testing.T) {
			worker := wpool.NewWorker(ctx, test.cfg)
			var wg sync.WaitGroup
			for i := 0; i <= test.tasksCnt; i++ {
				wg.Add(1)
				assert.NoError(tt, worker.Do(&wpool.Task{Callback: func(tCtx context.Context) {
					defer wg.Done()
					time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
					return
				}}))
			}
			wg.Wait()
			assert.NoError(tt, worker.Close())
		})
	}
}

func Test_Worker_DoAfterClose(t *testing.T) {
	worker := wpool.NewWorker(context.Background(), wpool.WorkerCfg{})
	worker.Close()
	assert.EqualError(t, worker.Do(&mocks.ITask{}), "add task to worker error: worker is closed")
}

func Test_Worker_ClosePipeline(t *testing.T) {
	pipeline := make(chan wpool.IWorker)
	wpool.NewWorker(context.Background(), wpool.WorkerCfg{Pipeline: pipeline, TasksChanBuff: 1})
	worker := <-pipeline
	close(pipeline)
	task := &mocks.ITask{}
	task.On("Do", mock.Anything).Once()
	worker.Do(task)
}
