package wpool_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/egnd/wpool"
	"github.com/rs/zerolog"
)

var logger = zerolog.Nop()

func benchPool(cfg wpool.PoolCfg, b *testing.B) {
	ctx := context.Background()
	pool := wpool.NewPool(cfg, func(num uint, pipeline chan wpool.IWorker) wpool.IWorker {
		return wpool.NewWorker(ctx, wpool.WorkerCfg{}, pipeline, &logger)
	}, &logger)
	defer pool.Close()

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		if err := pool.Add(&wpool.Task{Wg: &wg, Callback: func(tCtx context.Context, task *wpool.Task) error {
			select {
			case <-time.After(time.Millisecond):
				return nil
			case <-tCtx.Done():
				return &wpool.ErrTaskTimeout{task.GetName()}
			}
		}}); err != nil {
			b.Error(err)
			break
		}
		wg.Add(1)
	}
	wg.Wait()
}

func Benchmark_Pool_W10_T100(b *testing.B) {
	benchPool(wpool.PoolCfg{WorkersCnt: 10, TasksBufSize: 100}, b)
}

func Benchmark_Pool_W1_T100(b *testing.B) {
	benchPool(wpool.PoolCfg{WorkersCnt: 1, TasksBufSize: 100}, b)
}

func Benchmark_Pool_W10_T0(b *testing.B) {
	benchPool(wpool.PoolCfg{WorkersCnt: 10}, b)
}

func Benchmark_Pool_W1_T0(b *testing.B) {
	benchPool(wpool.PoolCfg{WorkersCnt: 1}, b)
}
