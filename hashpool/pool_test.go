package hashpool_test

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/egnd/go-pipeline"
	"github.com/egnd/go-pipeline/hashpool"
	"github.com/egnd/go-pipeline/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Pool(t *testing.T) {
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

			workers := []pipeline.Doer{}
			for i := 0; i < len(test.res); i++ {
				i := i
				worker := &mocks.Doer{}
				defer worker.AssertExpectations(tt)
				worker.On("Do", mock.Anything).After(time.Duration(rand.Intn(10)) * time.Millisecond).Run(func(args mock.Arguments) {
					task := args.Get(0).(pipeline.Task)
					res[i] = append(res[i], task.ID())
					task.Do()
				}).Return(nil).Maybe()
				worker.On("Close").Return(nil).Once()
				workers = append(workers, worker)
			}

			pipe := hashpool.NewPool(hashpool.DefaultHasher, workers...)

			var wg sync.WaitGroup
			for _, taskID := range test.tasks {
				wg.Add(1)
				task := &mocks.Task{}
				defer task.AssertExpectations(tt)
				task.On("ID").Return(taskID).Times(2)
				task.On("Do").After(time.Duration(rand.Intn(10)) * time.Millisecond).Run(func(_ mock.Arguments) { wg.Done() }).Once().Return(nil)
				assert.NoError(tt, pipe.Push(task))
			}

			wg.Wait()
			assert.NoError(tt, pipe.Close())
			assert.EqualValues(tt, test.res, res)
		})
	}
}

func Test_Pool_Errors(t *testing.T) {
	worker := &mocks.Doer{}
	worker.On("Close").Return(errors.New("worker close error")).Once()
	defer worker.AssertExpectations(t)

	pool := hashpool.NewPool(hashpool.DefaultHasher, worker)

	assert.EqualError(t, pool.Close(), "worker close error")
	assert.EqualError(t, pool.Close(), "pool close err: close of closed channel")
	assert.EqualError(t, pool.Push(nil), "pool push err: send on closed channel")
}

// func Test_Pool_Panic(t *testing.T) { @TODO:
// 	defer func() {
// 		if r := recover(); r != nil {
// 			log.Println("------", r)
// 			assert.EqualValues(t, "worker do error", r)
// 		}
// 	}()

// 	worker := &mocks.Doer{}
// 	worker.On("Do", mock.Anything).Return(errors.New("worker do error")).Once()
// 	// defer worker.AssertExpectations(t)

// 	pipe := hashpool.NewPool(hashpool.DefaultHasher, worker)

// 	task := &mocks.Task{}
// 	task.On("ID").Return("asdfg").Once()

// 	pipe.Push(task)
// }
