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
)

func Test_Pool(t *testing.T) {
	cases := []struct {
		pCfg     wpool.PoolCfg
		wCfg     wpool.WorkerCfg
		tasksCnt int
		taskErr  error
	}{
		{
			pCfg:     wpool.PoolCfg{TasksBufSize: 4, WorkersCnt: 3},
			tasksCnt: 10,
			wCfg:     wpool.WorkerCfg{TaskTTL: time.Millisecond},
		},
		{
			pCfg:     wpool.PoolCfg{WorkersCnt: 1},
			tasksCnt: 500,
		},
		{
			pCfg: wpool.PoolCfg{TasksBufSize: 20, WorkersCnt: 5},
		},
		{
			pCfg:     wpool.PoolCfg{WorkersCnt: 2},
			tasksCnt: 14,
		},
		{
			pCfg:     wpool.PoolCfg{TasksBufSize: 25, WorkersCnt: 1},
			tasksCnt: 13,
		},
	}
	ctx := context.Background()
	for k, test := range cases {
		t.Run(fmt.Sprint(k), func(tt *testing.T) {
			pool := wpool.NewPool(test.pCfg, func(num uint, pipeline chan wpool.IWorker) wpool.IWorker {
				test.wCfg.Pipeline = pipeline
				return wpool.NewWorker(ctx, test.wCfg)
			})

			var wg sync.WaitGroup
			for i := 0; i <= test.tasksCnt; i++ {
				wg.Add(1)
				go func() {
					assert.NoError(tt, pool.Add(&wpool.Task{Callback: func(tCtx context.Context) {
						defer wg.Done()
						time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
						return
					}}))
				}()
			}

			wg.Wait()

			assert.NoError(tt, pool.Close())
		})
	}
}

func Test_Pool_CloseError(t *testing.T) {
	pool := wpool.NewPool(wpool.PoolCfg{WorkersCnt: 10}, func(num uint, pipeline chan wpool.IWorker) wpool.IWorker {
		w := &mocks.IWorker{}
		w.On("Close").Return(errors.New("error"))
		return w
	})

	assert.Error(t, pool.Close())
	assert.Error(t, pool.Add(&mocks.ITask{}))
}
