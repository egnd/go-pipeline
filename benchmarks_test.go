package wpool_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/egnd/go-wpool"
)

func benchPool(cfg wpool.PoolCfg, b *testing.B) {
	ctx := context.Background()
	pool := wpool.NewPool(cfg, func(num uint, pipeline chan wpool.IWorker) wpool.IWorker {
		return wpool.NewWorker(ctx, wpool.WorkerCfg{
			Pipeline: pipeline,
			TaskTTL:  5 * time.Millisecond,
		})
	})
	defer pool.Close()

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		if err := pool.Add(&wpool.Task{Callback: func(ctx context.Context) {
			defer wg.Done()

			select {
			case <-time.After(time.Millisecond):
				return
			case <-ctx.Done():
				return
			}
		}}); err != nil {
			b.Error(err)
			break
		}
	}
	wg.Wait()
}

func benchWorker(cfg wpool.WorkerCfg, b *testing.B) {
	ctx := context.Background()
	worker := wpool.NewWorker(ctx, cfg)
	defer worker.Close()

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		if err := worker.Do(&wpool.Task{Callback: func(ctx context.Context) {
			defer wg.Done()

			select {
			case <-time.After(time.Millisecond):
				return
			case <-ctx.Done():
				return
			}
		}}); err != nil {
			b.Error(err)
			break
		}
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
