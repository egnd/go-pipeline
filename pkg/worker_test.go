package pkg_test

import (
	"testing"

	"github.com/egnd/wpool/mocks"
	"github.com/egnd/wpool/pkg"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_Worker(t *testing.T) {
	logger := zerolog.Nop()
	pipeline := make(chan pkg.IWorker)
	worker := pkg.NewWorker(pipeline, &logger)
	worker.Close()
	assert.EqualValues(t, pkg.ErrIsClosed(pkg.ErrIsClosed{EntityName: "worker"}), worker.Do(&mocks.ITask{}))
}
