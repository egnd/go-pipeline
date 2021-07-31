package wpool_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/egnd/wpool"
	"github.com/egnd/wpool/mocks"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func Test_Worker(t *testing.T) {
	cases := []struct {
		cfg      wpool.WorkerCfg
		tasksCnt int
		err      error
	}{
		{
			cfg: wpool.WorkerCfg{
				TasksChanBuff: 10,
				TaskTTL:       15 * time.Millisecond,
			},
			tasksCnt: 21,
		},
	}
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	ctx := context.Background()
	for k, test := range cases {
		t.Run(fmt.Sprint(k), func(tt *testing.T) {
			wLog := logger.With().Int("worker", k).Logger()
			worker := wpool.NewWorker(ctx, test.cfg, &wLog)
			defer worker.Close()

			for i := 0; i <= test.tasksCnt; i++ {
				if err := worker.Do(&wpool.Task{Name: fmt.Sprint(i), Callback: func(tCtx context.Context, task *wpool.Task) error {
					time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
					return test.err
				}}); err != nil {
					break
				}
			}
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		})
	}
}

func Test_Worker_DoAfterClose(t *testing.T) {
	logger := zerolog.Nop()
	worker := wpool.NewWorker(context.Background(), wpool.WorkerCfg{
		Pipeline: make(chan wpool.IWorker),
	}, &logger)
	worker.Close()
	assert.EqualValues(t, wpool.ErrIsClosed(wpool.ErrIsClosed{EntityName: "worker"}), worker.Do(&mocks.ITask{}))
}
