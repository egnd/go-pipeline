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
	TasksChanBuff uint
	TaskTTL       time.Duration
	Pipeline      chan<- IWorker
}

// Worker is a struct for handling tasks.
type Worker struct {
	closed bool
	cfg    WorkerCfg
	mx     sync.Mutex
	tasks  chan ITask
	logger *zerolog.Logger
}

// NewWorker is a factory method for creating of new workers.
func NewWorker(ctx context.Context, cfg WorkerCfg,
	logger *zerolog.Logger,
) IWorker {
	w := &Worker{ //nolint:exhaustivestruct
		tasks:  make(chan ITask, cfg.TasksChanBuff),
		logger: logger,
		cfg:    cfg,
	}

	w.logger.Debug().Msg("spawned")

	go w.run(ctx)

	return w
}

func (w *Worker) run(ctx context.Context) {
	for {
		w.notifyPipeline()

		for task := range w.tasks {
			w.logger.Debug().Str("task", task.GetName()).Msg("new task")

			if err := w.exec(ctx, task); err != nil {
				w.logger.Error().Str("task", task.GetName()).Err(err).Msg("do")
			}

			if w.cfg.Pipeline != nil {
				break
			}
		}
	}
}

func (w *Worker) notifyPipeline() {
	if w.cfg.Pipeline == nil {
		return
	}

	w.mx.Lock()
	defer w.mx.Unlock()

	if !w.closed {
		w.cfg.Pipeline <- w
	}
}

func (w *Worker) exec(ctx context.Context, task ITask) (tErr error) {
	defer func() {
		if panicMsg := recover(); panicMsg != nil {
			tErr = ErrPanic{panicMsg}
		}
	}()

	tCtx := ctx

	if w.cfg.TaskTTL > 0 {
		var tCtxCancel context.CancelFunc
		tCtx, tCtxCancel = context.WithTimeout(ctx, w.cfg.TaskTTL)

		defer tCtxCancel()
	}

	if tErr = task.Do(tCtx); tErr != nil {
		tErr = ErrWrapper{Msg: "task execution error", Err: tErr}
	}

	return
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
