# wpool

Golang package for making a pool of workers.

[![Go Reference](https://pkg.go.dev/badge/github.com/egnd/wpool.svg)](https://pkg.go.dev/github.com/egnd/wpool)
[![Go Report Card](https://goreportcard.com/badge/github.com/egnd/wpool)](https://goreportcard.com/report/github.com/egnd/wpool)
[![Coverage](http://gocover.io/_badge/github.com/egnd/wpool)](http://gocover.io/github.com/egnd/wpool)
[![Pipeline](https://github.com/egnd/wpool/actions/workflows/pipeline.yml/badge.svg)](https://github.com/egnd/wpool/actions?query=workflow%3APipeline)

### Pool example:
```golang
logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
ctx := context.Background()

// create pool and define worker's factory
pool := wpool.NewPool(wpool.PoolCfg{
    WorkersCnt:   3,
    TasksBufSize: 10,
}, func(num uint, pipeline chan wpool.IWorker) wpool.IWorker {
    wLog := logger.With().Uint("worker", num).Logger()
    return wpool.NewWorker(ctx, wpool.WorkerCfg{
        Pipeline: pipeline,
        TaskTTL:  300 * time.Millisecond,
    }, &wLog)
}, &logger)
defer pool.Close()

// put some tasks to pool
var wg sync.WaitGroup
for i := 0; i < 20; i++ {
    if err := pool.Add(&wpool.Task{Wg: &wg, Name: fmt.Sprint(i),
        Callback: func(tCtx context.Context, task *wpool.Task) error {
            // do something here
            logger.Info().Str("task", task.GetName()).Msg("do task")
            return nil
        },
    }); err != nil {
        logger.Error().Err(err).Msg("putting task to pool")
        break
    }
    wg.Add(1)
}
// wait for tasks to be completed
wg.Wait()
```

### Worker example:
```golang
logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
ctx := context.Background()

// create worker
worker := wpool.NewWorker(ctx, wpool.WorkerCfg{
    TasksChanBuff: 10,
    TaskTTL:       300 * time.Duration
}, &logger)
defer worker.Close()

// put some tasks to worker
var wg sync.WaitGroup
for i := 0; i < 20; i++ {
    if err := worker.Do(&wpool.Task{Wg: &wg, Name: fmt.Sprint(i),
        Callback: func(tCtx context.Context, task *wpool.Task) error {
            // do something here
            logger.Info().Str("task", task.GetName()).Msg("do task")
            return nil
        },
    }); err != nil {
        logger.Error().Err(err).Msg("putting task to worker")
        break
    }
    wg.Add(1)
}
// wait for tasks to be completed
wg.Wait()
```
