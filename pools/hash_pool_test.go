package pools_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/egnd/go-pipeline/assign"
	"github.com/egnd/go-pipeline/mocks"
	"github.com/egnd/go-pipeline/pools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_HashPool(t *testing.T) {
	cases := []struct {
		threadsCnt int
		tasksCnt   int
	}{
		{3, 100},
		{1, 20},
		{10, 200},
	}
	for k, test := range cases {
		t.Run(fmt.Sprint(k), func(tt *testing.T) {
			var wg sync.WaitGroup
			wg.Add(test.tasksCnt + 1)

			pipe := pools.NewHashPool(test.threadsCnt, 0, assign.Sticky)

			task := &mocks.Task{}
			defer task.AssertExpectations(tt)
			task.On("ID").Return("some-id").Times(test.tasksCnt)
			task.On("Do").After(time.Duration(rand.Intn(10)) * time.Millisecond).Run(func(_ mock.Arguments) { wg.Done() }).Times(test.tasksCnt).Return(nil)

			go func() {
				defer wg.Done()
				for i := 0; i < test.tasksCnt; i++ {
					assert.NoError(tt, pipe.Push(task))
				}
			}()

			wg.Wait()
			assert.NoError(tt, pipe.Close())
		})
	}
}

func Test_HashPool_Errors(t *testing.T) {
	pool := pools.NewHashPool(1, 0, assign.Sticky)

	assert.NoError(t, pool.Close())
	assert.EqualError(t, pool.Close(), "pool close err: close of closed channel")
	assert.EqualError(t, pool.Push(nil), "pool push err: send on closed channel")
}

func Test_HashPool_NoThreads(t *testing.T) {
	assert.PanicsWithValue(t, "HashPool requires at least 1 thread", func() {
		pools.NewHashPool(0, 0, assign.Sticky)
	})
}
