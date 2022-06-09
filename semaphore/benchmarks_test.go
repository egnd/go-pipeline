package semaphore_test

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/egnd/go-pipeline"
	"github.com/egnd/go-pipeline/semaphore"
)

func Benchmark_Semaphore(b *testing.B) {
	workCnt := []int{1, 10, 20}
	decoratorsCnt := []int{0, 1, 10}

	for _, wCnt := range workCnt {
		for _, dCnt := range decoratorsCnt {
			var wg sync.WaitGroup

			task := defTask{&wg}

			decorators := make([]pipeline.DoerDecorator, 0, dCnt)
			for i := 0; i < dCnt; i++ {
				decorators = append(decorators, defDecorator)
			}

			b.Run(fmt.Sprintf("w%d_d%d", wCnt, dCnt), func(bb *testing.B) {
				pipe := semaphore.NewSemaphore(wCnt, decorators...)

				wg.Add(bb.N)
				for k := 0; k < bb.N; k++ {
					if err := pipe.Push(&task); err != nil {
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

func defDecorator(next pipeline.Tasker) pipeline.Tasker {
	return func(task pipeline.Task) error {
		if task.ID() == "asfdsagsgsf" {
			log.Println("asdasdasd")
		}

		return next(task)
	}
}

type defTask struct {
	wg *sync.WaitGroup
}

func (t *defTask) ID() string { return "default task" }

func (t *defTask) Do() error {
	defer t.wg.Done()
	return nil
}
