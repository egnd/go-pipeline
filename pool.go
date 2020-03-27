package wpool

import (
	"sync"
)

type PoolInterface interface {
	RegisterWorker(worker WorkerInterface) PoolInterface
	AddJob(job JobInterface) PoolInterface
	Start() PoolInterface
	Wait() PoolInterface
	Stop() PoolInterface
	WorkersChannel() chan WorkerInterface
	WorkersGroup() *sync.WaitGroup
}

// @TODO: логирование
type Pool struct {
	PoolInterface
	workers     []WorkerInterface
	jobs        []JobInterface
	workersChan chan WorkerInterface
	stop        chan bool
	wg          sync.WaitGroup
	once        sync.Once
}

func (this *Pool) RegisterWorker(worker WorkerInterface) PoolInterface {
	this.workers = append(this.workers, worker)
	return this
}

func (this *Pool) AddJob(job JobInterface) PoolInterface {
	this.jobs = append(this.jobs, job)
	this.wg.Add(1)
	return this
}

func (this *Pool) Start() PoolInterface {
	go func() {
		var job JobInterface
		for {
			select {
			case <-this.stop:
				break
			case worker := <-this.workersChan:
				if len(this.jobs) > 0 {
					job, this.jobs = this.jobs[0], this.jobs[1:] // @TODO: облегчить, возможно переделать очередь в стек
					worker.Execute(job)
				}
			}
		}
	}()
	return this
}

func (this *Pool) Wait() PoolInterface {
	this.wg.Wait()
	return this
}

func (this *Pool) Stop() PoolInterface {
	for _, worker := range this.workers {
		worker.Stop()
	}
	this.stop <- true
	this.once.Do(func() {
		close(this.workersChan)
		close(this.stop)
	})
	return this
}

func (this *Pool) WorkersChannel() chan WorkerInterface {
	return this.workersChan
}

func (this *Pool) WorkersGroup() *sync.WaitGroup {
	return &this.wg
}

func NewPool() PoolInterface {
	return &Pool{
		workersChan: make(chan WorkerInterface),
		stop:        make(chan bool),
	}
}
