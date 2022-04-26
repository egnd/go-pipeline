package wpool_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/egnd/go-wpool/v2"
	"github.com/egnd/go-wpool/v2/interfaces"
	"github.com/rs/zerolog"
)

type SomeTask struct {
	wg    *sync.WaitGroup
	delay time.Duration
	id    string
}

func (t *SomeTask) GetID() string { return t.id }

func (t *SomeTask) Do() {
	defer t.wg.Done()
	if t.delay > 0 {
		time.Sleep(t.delay)
	}
}

func benchPipelinePool(workerCnt int, delay time.Duration, b *testing.B) {
	pipeline := make(chan interfaces.Worker)
	logger := wpool.NewZerologAdapter(zerolog.Nop())

	pool := wpool.NewPipelinePool(pipeline, logger)
	defer pool.Close()

	for i := 0; i < workerCnt; i++ {
		pool.AddWorker(wpool.NewPipelineWorker(pipeline))
	}

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)

		if err := pool.AddTask(
			&SomeTask{&wg, delay, "task" + fmt.Sprint(i)},
		); err != nil {
			b.Error(err)
			break
		}
	}
	wg.Wait()
}

func benchStickyPool(workerCnt int, workerSize int, delay time.Duration, b *testing.B) {
	logger := wpool.NewZerologAdapter(zerolog.Nop())

	pool := wpool.NewStickyPool(logger)
	defer pool.Close()

	for i := 0; i < workerCnt; i++ {
		pool.AddWorker(wpool.NewWorker(workerSize))
	}

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)

		if err := pool.AddTask(
			&SomeTask{&wg, delay, "task" + fmt.Sprint(i)},
		); err != nil {
			b.Error(err)
			break
		}
	}
	wg.Wait()
}

func Benchmark_PPool_Threads1(b *testing.B) {
	benchPipelinePool(1, 10*time.Millisecond, b)
}

func Benchmark_PPool_Threads10(b *testing.B) {
	benchPipelinePool(10, 10*time.Millisecond, b)
}

func Benchmark_SPool_Threads1_WSize1(b *testing.B) {
	benchStickyPool(1, 1, 10*time.Millisecond, b)
}

func Benchmark_SPool_Threads1_WSize10(b *testing.B) {
	benchStickyPool(1, 10, 10*time.Millisecond, b)
}

func Benchmark_SPool_Threads10_WSize1(b *testing.B) {
	benchStickyPool(10, 1, 10*time.Millisecond, b)
}

func Benchmark_SPool_Threads10_WSize100(b *testing.B) {
	benchStickyPool(10, 10, 10*time.Millisecond, b)
}
