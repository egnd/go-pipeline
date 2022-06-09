// Package hashpool contains pool and worker structs
package hashpool

import (
	"fmt"

	"github.com/egnd/go-pipeline"
)

// Pool is a pool of "sticky" workers.
type Pool struct {
	tasks   chan pipeline.Task
	workers []pipeline.Doer
}

// NewPool creates pool of "sticky" workers.
func NewPool(hasher Hasher, workers ...pipeline.Doer) *Pool {
	pool := &Pool{
		workers: workers,
		tasks:   make(chan pipeline.Task),
	}

	go func() {
		for task := range pool.tasks {
			worker := pool.workers[hasher(task.ID(), uint64(len(pool.workers)))]

			if err := worker.Do(task); err != nil {
				panic(err)
			}
		}
	}()

	return pool
}

// Push is putting task into pool.
func (p *Pool) Push(task pipeline.Task) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("pool push err: %v", r)
		}
	}()

	p.tasks <- task

	return
}

// Close is stopping pool and workers.
func (p *Pool) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("pool close err: %v", r)
		}
	}()

	close(p.tasks)

	for _, worker := range p.workers {
		if err = worker.Close(); err != nil {
			return
		}
	}

	return
}
