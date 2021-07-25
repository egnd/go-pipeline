package pkg_test

import (
	"sync"
	"testing"
	"time"

	"github.com/egnd/wpool/pkg"
	"github.com/rs/zerolog"
)

func benchPool(cfg pkg.PoolCfg, b *testing.B) {
	// b.Log("tasks", b.N, "buffer", cfg.TasksBufSize, "workers", cfg.WorkersCnt)
	logger := zerolog.Nop()
	// log.Print("start pool")
	pool := pkg.NewPool(cfg, func(num uint, pipeline chan pkg.IWorker) pkg.IWorker {
		return pkg.NewWorker(pipeline, &logger)
	}, &logger)

	// log.Print("add tasks")
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		if err := pool.Add(&pkg.Task{Title: "", Wg: &wg, Callback: func() error {
			// log.Print("run task")
			time.Sleep(time.Millisecond)
			return nil
		}}); err != nil {
			b.Error(err)
			break
		}
		wg.Add(1)
	}
	// log.Print("waiting tasks")
	wg.Wait()
	// log.Print("close pool")
	pool.Close()
}

func Benchmark_Pool_W10_T100(b *testing.B) {
	benchPool(pkg.PoolCfg{WorkersCnt: 10, TasksBufSize: 100}, b)
}

func Benchmark_Pool_W1_T100(b *testing.B) {
	benchPool(pkg.PoolCfg{WorkersCnt: 1, TasksBufSize: 100}, b)
}

func Benchmark_Pool_W10_T0(b *testing.B) {
	benchPool(pkg.PoolCfg{WorkersCnt: 10}, b)
}

func Benchmark_Pool_W1_T0(b *testing.B) {
	benchPool(pkg.PoolCfg{WorkersCnt: 1}, b)
}
