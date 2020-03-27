package wpool

import (
	"sync"
)

type WorkerInterface interface {
	Name() string
	Execute(job JobInterface) WorkerInterface
	Stop() WorkerInterface
}

type Worker struct {
	WorkerInterface
	name string
	wg   *sync.WaitGroup
	stop chan bool
	once sync.Once
}

func (this *Worker) Name() string {
	return this.name
}

func (this *Worker) Execute(job JobInterface) WorkerInterface {
	go func() {
		defer this.wg.Done()
		if job != nil {
			job.Run()
		}
	}()
	return this
}

func (this *Worker) Stop() WorkerInterface {
	this.stop <- true
	this.once.Do(func() {
		close(this.stop)
	})
	return this
}

func NewWorker(name string, pool PoolInterface) (worker WorkerInterface) {
	done := make(chan bool)
	worker = &Worker{
		name: name,
		wg:   pool.WorkersGroup(),
		stop: done,
	}
	go func() {
		for {
			select {
			case <-done:
				break
			default:
				pool.WorkersChannel() <- worker
			}
		}
	}()
	return
}
