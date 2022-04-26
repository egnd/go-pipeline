package wpool_test

import (
	"testing"

	"github.com/egnd/go-wpool/v2"
	"github.com/rs/zerolog"
)

func Test_ZerologAdapter(t *testing.T) {
	logger := wpool.NewZerologAdapter(zerolog.Nop())
	logger.Errorf(nil, "test error")
	logger.Infof("test info")
}
