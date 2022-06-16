# go-pipeline

[![Go Reference](https://pkg.go.dev/badge/github.com/egnd/go-pipeline.svg)](https://pkg.go.dev/github.com/egnd/go-pipeline)
[![Go Report Card](https://goreportcard.com/badge/github.com/egnd/go-pipeline)](https://goreportcard.com/report/github.com/egnd/go-pipeline)
[![Coverage](https://gocover.io/_badge/github.com/egnd/go-pipeline?k1)](https://gocover.io/github.com/egnd/go-pipeline)
[![Pipeline](https://github.com/egnd/go-pipeline/actions/workflows/pipeline.yml/badge.svg)](https://github.com/egnd/go-pipeline/actions?query=workflow%3APipeline)

Golang package for parallel execution of tasks.

### BusPool:
Common Pool of Workers. The Task is taken into work by the first released Worker.
```golang
package main

import (
	"sync"

	"github.com/egnd/go-pipeline/decorators"
	"github.com/egnd/go-pipeline/pools"
	"github.com/egnd/go-pipeline/tasks"
	"github.com/rs/zerolog"
)

func main() {
	// create pool
	pipe := pools.NewBusPool(
		2,  // set parallel threads count
		10, // set tasks queue size
		// add some task decorators:
		decorators.LogErrorZero(zerolog.Nop()), // log tasks errors
		decorators.CatchPanic,                  // convert tasks panics to errors
	)

	// start producing tasks to pool
	var wg sync.WaitGroup
	pipe.Push(tasks.NewFunc("testid", func() error {
		defer wg.Done()
		return nil
	}))
	pipe.Push(tasks.NewFunc("testid", func() error {
		defer wg.Done()
		return nil
	}))

	wg.Wait()

	// close pool
	if err := pipe.Close(); err != nil {
		panic(err)
	}
}
```

### HashPool:
Worker pool, which allows you to change the strategy for assigning Tasks to Workers.
```golang
package main

import (
	"sync"

	"github.com/egnd/go-pipeline/assign"
	"github.com/egnd/go-pipeline/decorators"
	"github.com/egnd/go-pipeline/pools"
	"github.com/egnd/go-pipeline/tasks"
	"github.com/rs/zerolog"
)

func main() {
	// create pool
	pipe := pools.NewHashPool(
		2,             // set parallel threads count
		10,            // set tasks queue size
		assign.Sticky, // choose tasks to workers assignment method
		// add some task decorators:
		decorators.LogErrorZero(zerolog.Nop()), // log tasks errors
		decorators.CatchPanic,                  // convert tasks panics to errors
	)

	// start producing tasks to pool
	var wg sync.WaitGroup
	pipe.Push(tasks.NewFunc("testid", func() error {
		defer wg.Done()
		return nil
	}))
	pipe.Push(tasks.NewFunc("testid", func() error {
		defer wg.Done()
		return nil
	}))

	wg.Wait()

	// close pool
	if err := pipe.Close(); err != nil {
		panic(err)
	}
}
```

### Semaphore:
Primitive for limiting the number of threads for the Tasks parallel execution.
```golang
package main

import (
	"sync"

	"github.com/egnd/go-pipeline/decorators"
	"github.com/egnd/go-pipeline/pools"
	"github.com/egnd/go-pipeline/tasks"
	"github.com/rs/zerolog"
)

func main() {
	// create pool
	pipe := pools.NewSemaphore(2, // set parallel threads count
		// add some task decorators:
		decorators.LogErrorZero(zerolog.Nop()), // log tasks errors
		decorators.CatchPanic,                  // convert tasks panics to errors
	)

	// start producing tasks to pool
	var wg sync.WaitGroup
	pipe.Push(tasks.NewFunc("testid", func() error {
		defer wg.Done()
		return nil
	}))
	pipe.Push(tasks.NewFunc("testid", func() error {
		defer wg.Done()
		return nil
	}))

	wg.Wait()

	// close pool
	if err := pipe.Close(); err != nil {
		panic(err)
	}
}
```
