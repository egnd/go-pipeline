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

func Test_StickyPool(t *testing.T) {
	cases := []struct {
		tasks []string
		res   [][]string
	}{
		{
			tasks: []string{
				"asdfg", "asdfg", "asdfg",
			},
			res: [][]string{
				{}, {}, {"asdfg", "asdfg", "asdfg"},
			},
		},
		{
			tasks: []string{
				"asdfg", "asddsafdsgfg", "asdfg",
				"asdf234235g", "a====sdfg", "asds-df=sfg",
				"asdfa8dsf7sag", "asdfg", "f0gdsf00sasdfg",
			},
			res: [][]string{
				{"asdf234235g", "a====sdfg"},
				{"asddsafdsgfg", "asds-df=sfg"},
				{"asdfg", "asdfg", "asdfa8dsf7sag", "asdfg", "f0gdsf00sasdfg"},
			},
		},
		{
			tasks: []string{
				"asdfg", "asddsafdsgfg", "asdfg",
				"asdf234235g", "a====sdfg", "asds-df=sfg",
				"asdfa8dsf7sag", "asdfg", "f0gdsf00sasdfg",
			},
			res: [][]string{
				{
					"asdfg", "asddsafdsgfg", "asdfg", "asdf234235g", "a====sdfg",
					"asds-df=sfg", "asdfa8dsf7sag", "asdfg", "f0gdsf00sasdfg",
				},
			},
		},
	}
	for k, test := range cases {
		t.Run(fmt.Sprint(k), func(tt *testing.T) {
			var res [][]string
			for range test.res {
				res = append(res, []string{})
			}

			pool := wpool.NewStickyPool(nil)

			for i := 0; i < len(test.res); i++ {
				i := i
				worker := &interfaces.MockWorker{}
				defer worker.AssertExpectations(tt)
				worker.On("Do", mock.Anything).After(time.Duration(rand.Intn(10)) * time.Millisecond).Run(func(args mock.Arguments) {
					task := args.Get(0).(interfaces.Task)
					res[i] = append(res[i], task.GetID())
					task.Do()
				}).Return(nil).Maybe()
				worker.On("Close").Return(nil).Once()
				pool.AddWorker(worker)
			}

			var wg sync.WaitGroup
			for _, taskID := range test.tasks {
				wg.Add(1)
				task := &interfaces.MockTask{}
				defer task.AssertExpectations(tt)
				task.On("GetID").Return(taskID).Times(2)
				task.On("Do").After(time.Duration(rand.Intn(10)) * time.Millisecond).Run(func(_ mock.Arguments) { wg.Done() }).Once()
				assert.NoError(tt, pool.AddTask(task))
			}

			wg.Wait()
			assert.NoError(tt, pool.Close())
			assert.EqualValues(tt, test.res, res)
		})
	}
}

func Test_StickyPool_Close_Error(t *testing.T) {
	pool := wpool.NewStickyPool(nil)

	worker := &interfaces.MockWorker{}
	worker.On("Close").Return(errors.New("error")).Once()
	defer worker.AssertExpectations(t)
	pool.AddWorker(worker)

	assert.EqualError(t, pool.Close(), "error")
	assert.EqualError(t, pool.AddTask(nil), "pool is closed")
}

func Test_StickyPool_Do_Error(t *testing.T) {
	logger := &interfaces.MockLogger{}
	logger.On("Errorf", mock.Anything, `worker #%d doing task "%s"`, uint64(0), "asdfggg")

	pool := wpool.NewStickyPool(logger)

	worker := &interfaces.MockWorker{}
	worker.On("Do", mock.Anything).Return(errors.New("error")).Once()
	worker.On("Close").Return(nil).Once()
	defer worker.AssertExpectations(t)
	pool.AddWorker(worker)

	task := &interfaces.MockTask{}
	task.On("GetID").Return("asdfggg").Times(2)
	defer task.AssertExpectations(t)
	pool.AddTask(task)

	time.Sleep(50 * time.Millisecond)
	assert.NoError(t, pool.Close())
}
