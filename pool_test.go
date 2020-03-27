package wpool

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	log.Print("--------------------- Testing pool")
	runPool(3, 10, true, time.Duration(rand.Intn(300))*time.Millisecond)
}

// func BenchmarkPoolSingleWorker(b *testing.B) {
// 	runPool(1, b.N, false, 100*time.Millisecond)
// }

// func BenchmarkPool10Workers(b *testing.B) {
// 	runPool(10, b.N, false, 100*time.Millisecond)
// }

// func BenchmarkPool100Workers(b *testing.B) {
// 	runPool(100, b.N, false, 100*time.Millisecond)
// }

func runPool(workersCnt int, jobsCnt int, debug bool, jobDelay time.Duration) {
	if debug {
		log.Print("start pool")
	}
	pool := NewPool().Start()
	defer func() {
		if debug {
			log.Print("stop pool")
		}
		pool.Stop()
	}()
	for i := 1; i <= workersCnt; i++ {
		if debug {
			log.Printf("add worker #%d", i)
		}
		worker := NewWorker(fmt.Sprintf("worker-%d", i), pool)
		pool.RegisterWorker(worker)
	}
	for i := 1; i <= jobsCnt; i++ {
		if debug {
			log.Printf("add job #%d", i)
		}
		pool.AddJob(NewJob(fmt.Sprintf("job-%d", i), func(job JobInterface) (err error) {
			if debug {
				log.Print(job.Name() + " in progress")
			}
			time.Sleep(jobDelay)
			return
		}))
	}
	pool.Wait()
}
