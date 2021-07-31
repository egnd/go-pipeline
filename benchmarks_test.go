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
		return wpool.NewWorker(ctx, wpool.WorkerCfg{
			Pipeline: pipeline,
			TaskTTL:  5 * time.Millisecond,
		}, &logger)
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

func benchWorker(cfg wpool.WorkerCfg, b *testing.B) {
	ctx := context.Background()
	worker := wpool.NewWorker(ctx, cfg, &logger)
	defer worker.Close()

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		if err := worker.Do(&wpool.Task{Wg: &wg, Callback: func(tCtx context.Context, task *wpool.Task) error {
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

func Benchmark_Pool_Workers10_Buff100(b *testing.B) {
	benchPool(wpool.PoolCfg{WorkersCnt: 10, TasksBufSize: 100}, b)
}

func Benchmark_Pool_Workers1_Buff100(b *testing.B) {
	benchPool(wpool.PoolCfg{WorkersCnt: 1, TasksBufSize: 100}, b)
}

func Benchmark_Pool_Workers10_Buff0(b *testing.B) {
	benchPool(wpool.PoolCfg{WorkersCnt: 10}, b)
}

func Benchmark_Pool_Workers1_Buff0(b *testing.B) {
	benchPool(wpool.PoolCfg{WorkersCnt: 1}, b)
}

func Benchmark_Worker_Buf0(b *testing.B) {
	benchWorker(wpool.WorkerCfg{TasksChanBuff: 0}, b)
}

func Benchmark_Worker_Buf1(b *testing.B) {
	benchWorker(wpool.WorkerCfg{TasksChanBuff: 1}, b)
}

func Benchmark_Worker_Buf10(b *testing.B) {
	benchWorker(wpool.WorkerCfg{TasksChanBuff: 10}, b)
}
