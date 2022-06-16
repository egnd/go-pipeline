package pipeline_test

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/egnd/go-pipeline"
	"github.com/egnd/go-pipeline/assign"
	"github.com/egnd/go-pipeline/pools"
	"github.com/egnd/go-pipeline/tasks"
)

func Benchmark_BusPool(b *testing.B) {
	workCnt := []int{1, 10, 20}
	buffCnt := []int{0, 10}
	decoratorsCnt := []int{0, 1, 10}

	for _, wCnt := range workCnt {
		for _, bCnt := range buffCnt {
			for _, dCnt := range decoratorsCnt {
				var wg sync.WaitGroup

				task := tasks.NewFunc("testid", func() error {
					defer wg.Done()
					return nil
				})

				decorators := make([]pipeline.TaskDecorator, 0, dCnt)
				for i := 0; i < dCnt; i++ {
					decorators = append(decorators, defDecorator)
				}

				b.Run(fmt.Sprintf("w%d_b%d_d%d", wCnt, bCnt, dCnt), func(bb *testing.B) {
					pipe := pools.NewBusPool(wCnt, bCnt, decorators...)

					wg.Add(bb.N)
					for k := 0; k < bb.N; k++ {
						if err := pipe.Push(task); err != nil {
							bb.Error(err)
						}
					}
					wg.Wait()

					if err := pipe.Close(); err != nil {
						bb.Error(err)
					}
				})
			}
		}
	}
}

func Benchmark_HashPool(b *testing.B) {
	workCnt := []int{1, 10, 20}
	buffCnt := []int{0, 10}
	decoratorsCnt := []int{0, 1, 10}

	for _, wCnt := range workCnt {
		for _, bCnt := range buffCnt {
			for _, dCnt := range decoratorsCnt {
				var wg sync.WaitGroup

				task := tasks.NewFunc("testid", func() error {
					defer wg.Done()
					return nil
				})

				decorators := make([]pipeline.TaskDecorator, 0, dCnt)
				for i := 0; i < dCnt; i++ {
					decorators = append(decorators, defDecorator)
				}

				b.Run(fmt.Sprintf("w%d_b%d_d%d", wCnt, bCnt, dCnt), func(bb *testing.B) {
					pipe := pools.NewHashPool(wCnt, bCnt, assign.Random, decorators...)

					wg.Add(bb.N)
					for k := 0; k < bb.N; k++ {
						if err := pipe.Push(task); err != nil {
							bb.Error(err)
						}
					}
					wg.Wait()

					if err := pipe.Close(); err != nil {
						bb.Error(err)
					}
				})
			}
		}
	}
}

func Benchmark_Semaphore(b *testing.B) {
	workCnt := []int{1, 10, 20}
	decoratorsCnt := []int{0, 1, 10}

	for _, wCnt := range workCnt {
		for _, dCnt := range decoratorsCnt {
			var wg sync.WaitGroup

			task := tasks.NewFunc("testid", func() error {
				defer wg.Done()
				return nil
			})

			decorators := make([]pipeline.TaskDecorator, 0, dCnt)
			for i := 0; i < dCnt; i++ {
				decorators = append(decorators, defDecorator)
			}

			b.Run(fmt.Sprintf("w%d_d%d", wCnt, dCnt), func(bb *testing.B) {
				pipe := pools.NewSemaphore(wCnt, decorators...)

				wg.Add(bb.N)
				for k := 0; k < bb.N; k++ {
					if err := pipe.Push(task); err != nil {
						bb.Error(err)
					}
				}
				wg.Wait()

				if err := pipe.Close(); err != nil {
					bb.Error(err)
				}
			})
		}
	}
}

func defDecorator(next pipeline.TaskExecutor) pipeline.TaskExecutor {
	return func(task pipeline.Task) error {
		if task.ID() == "asfdsagsgsf" {
			log.Println("asdasdasd")
		}

		return next(task)
	}
}
