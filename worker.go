package wpool

import (
	"sync"
)

type WorkerInterface interface {
	Name() string
	Execute(job JobInterface) WorkerInterface
	Stop() WorkerInterface
	PrependJob(callback WorkerCallback) WorkerInterface
	AppendJob(callback WorkerCallback) WorkerInterface
}

type WorkerCallback func(job JobInterface) error

type Worker struct {
	WorkerInterface
	name       string
	wg         *sync.WaitGroup
	stop       chan bool
	once       sync.Once
	prependJob WorkerCallback
	appendJob  WorkerCallback
}

func (this *Worker) Name() string {
	return this.name
}

func (this *Worker) Execute(job JobInterface) WorkerInterface {
	go func() {
		defer this.wg.Done()
		if job != nil {
			var err error
			if this.prependJob != nil {
				err = this.prependJob(job)
			}
			if err == nil {
				job.Run()
				if this.appendJob != nil {
					this.appendJob(job)
				}
			}
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

func (this *Worker) PrependJob(callback WorkerCallback) WorkerInterface {
	this.prependJob = callback
	return this
}

func (this *Worker) AppendJob(callback WorkerCallback) WorkerInterface {
	this.appendJob = callback
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
