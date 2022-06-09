package hashpool_test

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/egnd/go-pipeline"
	"github.com/egnd/go-pipeline/hashpool"
)

func Benchmark_Hashpool(b *testing.B) {
	workCnt := []int{1, 10, 20}
	buffCnt := []int{0, 10}
	decoratorsCnt := []int{0, 1, 10}

	for _, wCnt := range workCnt {
		for _, bCnt := range buffCnt {
			for _, dCnt := range decoratorsCnt {
				var wg sync.WaitGroup

				decorators := make([]pipeline.DoerDecorator, 0, dCnt)
				for i := 0; i < dCnt; i++ {
					decorators = append(decorators, defDecorator)
				}

				b.Run(fmt.Sprintf("w%d_b%d_d%d", wCnt, bCnt, dCnt), func(bb *testing.B) {
					workers := make([]pipeline.Doer, 0, wCnt)
					for i := 0; i < wCnt; i++ {
						workers = append(workers, hashpool.NewWorker(bCnt, decorators...))
					}

					pipe := hashpool.NewPool(hashpool.DefaultHasher, workers...)

					wg.Add(bb.N)
					for k := 0; k < bb.N; k++ {
						if err := pipe.Push(&defTask{k, &wg}); err != nil {
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

func defDecorator(next pipeline.Tasker) pipeline.Tasker {
	return func(task pipeline.Task) error {
		if task.ID() == "asfdsagsgsf" {
			log.Println("asdasdasd")
		}

		return next(task)
	}
}

type defTask struct {
	num int
	wg  *sync.WaitGroup
}

func (t *defTask) ID() string { return fmt.Sprintf("task#%d", t.num) }

func (t *defTask) Do() error {
	defer t.wg.Done()
	return nil
}
