package wpool_test

import (
	"context"
	"testing"

	"github.com/egnd/wpool"
	"github.com/egnd/wpool/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_Worker(t *testing.T) {
	logger := zerolog.Nop()
	pipeline := make(chan wpool.IWorker)
	worker := wpool.NewWorker(context.Background(), wpool.WorkerCfg{}, pipeline, &logger)
	worker.Close()
	assert.EqualValues(t, wpool.ErrIsClosed(wpool.ErrIsClosed{EntityName: "worker"}), worker.Do(&mocks.ITask{}))
}
