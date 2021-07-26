package wpool

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// WorkerFactory is a factory method to pass into pool of workers.
type WorkerFactory func(uint, chan IWorker) IWorker

// IWorker is a worker interface.
type IWorker interface {
	Do(ITask) error
	Close()
}

// WorkerCfg is a config for Worker.
type WorkerCfg struct {
	Timeout time.Duration
}

// Worker is a struct for handling tasks.
type Worker struct {
	closed bool
	mx     sync.Mutex
	tasks  chan ITask
	// @TODO: logger interface
	logger *zerolog.Logger
}

// NewWorker is a factory method for creating of new workers.
func NewWorker(ctx context.Context,
	cfg WorkerCfg, pipeline chan<- IWorker, logger *zerolog.Logger,
) IWorker {
	w := &Worker{ //nolint:exhaustivestruct
		tasks:  make(chan ITask),
		logger: logger,
	}

	w.logger.Debug().Msg("spawned")

	go func() {
		for {
			w.mx.Lock()

			if !w.closed {
				pipeline <- w
			}

			w.mx.Unlock()

			for task := range w.tasks {
				w.logger.Debug().Str("task", task.GetName()).Msg("new task")

				err := func() (tErr error) {
					defer func() {
						if panicMsg := recover(); panicMsg != nil {
							tErr = ErrPanic{panicMsg}
						}
					}()

					tCtx := ctx

					if cfg.Timeout > 0 {
						var tCtxCancel context.CancelFunc
						tCtx, tCtxCancel = context.WithTimeout(ctx, cfg.Timeout)

						defer tCtxCancel()
					}

					if tErr = task.Do(tCtx); tErr != nil {
						tErr = ErrWrapper{Msg: "task execution error", Err: tErr}
					}

					return
				}()
				if err != nil {
					w.logger.Error().Str("task", task.GetName()).Err(err).Msg("do")
				}

				break
			}
		}
	}()

	return w
}

// Do is method for putting task to worker.
func (w *Worker) Do(task ITask) error {
	w.mx.Lock()

	defer w.mx.Unlock()

	if w.closed {
		return ErrIsClosed{"worker"}
	}

	w.tasks <- task

	return nil
}

// Close is a method for worker stopping.
func (w *Worker) Close() {
	w.logger.Debug().Msg("close")

	w.mx.Lock()
	defer w.mx.Unlock()

	w.closed = true

	close(w.tasks)
}
