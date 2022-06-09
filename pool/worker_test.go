package pool_test

import (
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

func Test_Worker_Do(t *testing.T) {
	cases := []struct {
		buffSize int
		tasksCnt int
	}{
		{
			buffSize: 10,
			tasksCnt: 21,
		},
		{
			buffSize: 20,
			tasksCnt: 10,
		},
	}
	for k, test := range cases {
		t.Run(fmt.Sprint(k), func(tt *testing.T) {
			notifier := make(chan pipeline.Doer)
			pool.NewWorker(notifier)

			var wg sync.WaitGroup
			for i := 0; i <= test.tasksCnt; i++ {
				wg.Add(1)

				task := &mocks.Task{}
				defer task.AssertExpectations(tt)
				task.On("Do").Once().After(time.Duration(rand.Intn(10)) * time.Millisecond).Run(func(_ mock.Arguments) { wg.Done() }).Return(nil)

				assert.NoError(tt, (<-notifier).Do(task))
			}

			wg.Wait()
			assert.NoError(tt, (<-notifier).Close())
		})
	}
}

func Test_Worker_Close_Error(t *testing.T) {
	notifier := make(chan pipeline.Doer)
	pool.NewWorker(notifier)
	worker := <-notifier
	worker.Close()
	assert.EqualError(t, worker.Do(nil), "worker do err: send on closed channel")
}

func Test_WorkerPipeline_Pipeline_Close_Error(t *testing.T) {
	notifier := make(chan pipeline.Doer)
	close(notifier)
	w := pool.NewWorker(notifier)
	time.Sleep(50 * time.Millisecond)
	assert.EqualError(t, w.Close(), "worker close err: close of closed channel")
}
