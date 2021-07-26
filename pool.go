// Package wpool contains structs and functions for making a pool of workers.
package wpool

import (
	"sync"

	"github.com/rs/zerolog"
)

// IPool is pool of workers interface.
type IPool interface {
	Add(ITask) error
	Close()
}

// PoolCfg is pool config.
type PoolCfg struct {
	TasksBufSize uint
	WorkersCnt   uint
}

// Pool is struct for handlindg tasks with workers.
type Pool struct {
	closed   bool
	mx       sync.Mutex
	pipeline chan IWorker
	tasks    chan ITask
	workers  []IWorker
	// @TODO: logger interface
	logger *zerolog.Logger
}

// NewPool is a factory method for pool of workers.
func NewPool(cfg PoolCfg, worker WorkerFactory, logger *zerolog.Logger) IPool {
	p := &Pool{ //nolint:exhaustivestruct
		tasks:    make(chan ITask, cfg.TasksBufSize),
		pipeline: make(chan IWorker),
		workers:  make([]IWorker, 0, cfg.WorkersCnt),
		logger:   logger,
	}

	for i := uint(0); i < cfg.WorkersCnt; i++ {
		p.workers = append(p.workers, worker(i, p.pipeline))
	}

	go func() {
		for worker := range p.pipeline {
			for task := range p.tasks {
				if err := worker.Do(task); err != nil {
					p.logger.Error().Err(err).Msg("worker exec")
				}

				break //nolint:staticcheck
			}
		}
	}()

	return p
}

// Add is method for putting task into pool of workers.
func (p *Pool) Add(task ITask) error {
	p.mx.Lock()

	if p.closed {
		defer p.mx.Unlock()

		return ErrIsClosed{"pool"}
	}

	p.tasks <- task
	p.mx.Unlock()

	return nil
}

// Close is method for stopping for pool of workers.
func (p *Pool) Close() {
	p.mx.Lock()

	defer p.mx.Unlock()

	p.closed = true

	close(p.tasks)

	for _, worker := range p.workers {
		worker.Close()
	}

	close(p.pipeline)
}
