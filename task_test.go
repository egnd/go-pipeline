package wpool_test

import (
	"context"
	"sync"
	"testing"

	"github.com/egnd/go-wpool"
	"github.com/stretchr/testify/assert"
)

func Test_Task(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	assert.NoError(t, (&wpool.Task{Wg: &wg}).Do(context.Background()))
}
