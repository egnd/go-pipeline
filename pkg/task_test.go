package pkg_test

import (
	"sync"
	"testing"

	"github.com/egnd/wpool/pkg"
	"github.com/stretchr/testify/assert"
)

func Test_Task(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	assert.NoError(t, (&pkg.Task{Wg: &wg}).Do())
}
