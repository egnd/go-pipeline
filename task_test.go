package wpool_test

import (
	"context"
	"testing"

	"github.com/egnd/go-wpool"
)

func Test_Task(t *testing.T) {
	(&wpool.Task{Callback: func(ctx context.Context) { return }}).Do(context.Background())
	(&wpool.Task{}).Do(context.Background())
}
