package wpool_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/egnd/wpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type testCase struct {
	cfg       wpool.PoolCfg
	wCfg      wpool.WorkerCfg
	tasksCnt  int
	taskErr   error
	taskPanic string
}

func Test_Pool(t *testing.T) {
	cases := []testCase{
		{
			cfg:      wpool.PoolCfg{TasksBufSize: 4, WorkersCnt: 3},
			tasksCnt: 10,
			wCfg:     wpool.WorkerCfg{Timeout: time.Millisecond},
		},
		{
			cfg:      wpool.PoolCfg{WorkersCnt: 1},
			tasksCnt: 3000,
		},
		{
			cfg: wpool.PoolCfg{TasksBufSize: 20, WorkersCnt: 5},
		},
		{
			cfg:      wpool.PoolCfg{WorkersCnt: 2},
			tasksCnt: 14,
		},
		{
			cfg:      wpool.PoolCfg{TasksBufSize: 25, WorkersCnt: 1},
			tasksCnt: 13,
		},
		{
			cfg:      wpool.PoolCfg{WorkersCnt: 3},
			tasksCnt: 1,
			taskErr:  errors.New("error"),
		},
		{
			cfg:       wpool.PoolCfg{WorkersCnt: 3},
			tasksCnt:  1,
			taskPanic: "panic msg",
		},
	}
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	ctx := context.Background()
	for k, test := range cases {
		t.Run(fmt.Sprint(k), func(tt *testing.T) {
			pool := wpool.NewPool(test.cfg, func(num uint, pipeline chan wpool.IWorker) wpool.IWorker {
				wLog := logger.With().Uint("worker", num).Int("case", k).Logger()
				return wpool.NewWorker(ctx, test.wCfg, pipeline, &wLog)
			}, &logger)

			defer pool.Close()

			go func(tCase testCase) {
				for i := 0; i <= tCase.tasksCnt; i++ {
					if err := pool.Add(&wpool.Task{Name: fmt.Sprint(i), Callback: func(tCtx context.Context, task *wpool.Task) error {
						if len(tCase.taskPanic) > 0 {
							panic(tCase.taskPanic)
						}
						time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
						return tCase.taskErr
					}}); err != nil {
						break
					}
					time.Sleep(time.Duration(rand.Intn(2)) * time.Millisecond)
				}
			}(test)
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		})
	}
}
