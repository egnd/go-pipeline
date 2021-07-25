package pkg_test

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/egnd/wpool/pkg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type testCase struct {
	cfg       pkg.PoolCfg
	tasksCnt  int
	taskErr   error
	taskPanic string
}

func Test_Pool(t *testing.T) {
	cases := []testCase{
		{
			cfg:      pkg.PoolCfg{TasksBufSize: 4, WorkersCnt: 3},
			tasksCnt: 10,
		},
		{
			cfg:      pkg.PoolCfg{WorkersCnt: 1},
			tasksCnt: 3000,
		},
		{
			cfg: pkg.PoolCfg{TasksBufSize: 20, WorkersCnt: 5},
		},
		{
			cfg:      pkg.PoolCfg{WorkersCnt: 2},
			tasksCnt: 14,
		},
		{
			cfg:      pkg.PoolCfg{TasksBufSize: 25, WorkersCnt: 1},
			tasksCnt: 13,
		},
		{
			cfg:      pkg.PoolCfg{WorkersCnt: 3},
			tasksCnt: 1,
			taskErr:  errors.New("error"),
		},
		{
			cfg:       pkg.PoolCfg{WorkersCnt: 3},
			tasksCnt:  1,
			taskPanic: "panic msg",
		},
	}
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	for k, test := range cases {
		t.Run(fmt.Sprint(k), func(tt *testing.T) {
			pool := pkg.NewPool(test.cfg, func(num uint, pipeline chan pkg.IWorker) pkg.IWorker {
				wLog := logger.With().Uint("worker", num).Int("case", k).Logger()
				return pkg.NewWorker(pipeline, &wLog)
			}, &logger)
			go func(tCase testCase) {
				for i := 0; i <= tCase.tasksCnt; i++ {
					if err := pool.Add(&pkg.Task{Title: fmt.Sprint(i), Callback: func() error {
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
			pool.Close()
		})
	}
}
