package wpool_test

import (
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

func Test_WorkerPipeline_Do(t *testing.T) {
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
			pipeline := make(chan interfaces.Worker)
			wpool.NewPipelineWorker(pipeline)

			var wg sync.WaitGroup
			for i := 0; i <= test.tasksCnt; i++ {
				wg.Add(1)

				task := &interfaces.MockTask{}
				defer task.AssertExpectations(tt)
				task.On("Do").Once().After(time.Duration(rand.Intn(10)) * time.Millisecond).Run(func(_ mock.Arguments) { wg.Done() })

				assert.NoError(tt, (<-pipeline).Do(task))
			}

			wg.Wait()
			assert.NoError(tt, (<-pipeline).Close())
		})
	}
}

func Test_WorkerPipeline_Close_Error(t *testing.T) {
	pipeline := make(chan interfaces.Worker)
	wpool.NewPipelineWorker(pipeline)
	worker := <-pipeline
	worker.Close()
	assert.EqualError(t, worker.Do(nil), "worker is closed")
}

func Test_WorkerPipeline_Pipeline_Close_Error(t *testing.T) {
	pipeline := make(chan interfaces.Worker)
	close(pipeline)
	w := wpool.NewPipelineWorker(pipeline)
	time.Sleep(50 * time.Millisecond)
	assert.EqualError(t, w.Close(), "worker already closed")
}
