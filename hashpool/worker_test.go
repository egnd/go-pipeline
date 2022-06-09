package hashpool_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/egnd/go-pipeline/hashpool"
	"github.com/egnd/go-pipeline/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Worker(t *testing.T) {
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
			worker := hashpool.NewWorker(test.buffSize)

			var wg sync.WaitGroup
			for i := 0; i <= test.tasksCnt; i++ {
				wg.Add(1)

				task := &mocks.Task{}
				defer task.AssertExpectations(tt)
				task.On("Do").Once().After(time.Duration(rand.Intn(10)) * time.Millisecond).Run(func(_ mock.Arguments) { wg.Done() }).Return(nil)

				assert.NoError(tt, worker.Do(task))
			}

			wg.Wait()
			assert.NoError(tt, worker.Close())
		})
	}
}

func Test_Worker_Do_Error(t *testing.T) {
	worker := hashpool.NewWorker(0)
	assert.NoError(t, worker.Close())
	assert.EqualError(t, worker.Do(nil), "worker do err: send on closed channel")
}

func Test_Worker_Close_Error(t *testing.T) {
	worker := hashpool.NewWorker(0)
	assert.NoError(t, worker.Close())
	assert.EqualError(t, worker.Close(), "worker close err: close of closed channel")
}
