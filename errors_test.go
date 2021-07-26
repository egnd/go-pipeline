package wpool_test

import (
	"errors"
	"testing"

	"github.com/egnd/wpool"
	"github.com/stretchr/testify/assert"
)

func Test_ErrIsClosed(t *testing.T) {
	assert.EqualValues(t, (&wpool.ErrIsClosed{"pool"}).Error(), "pool is closed")
}

func Test_ErrPanic(t *testing.T) {
	assert.EqualValues(t, (&wpool.ErrPanic{"panic"}).Error(), "panic occurred: panic")
}

func Test_ErrWrapper(t *testing.T) {
	assert.EqualValues(t, (&wpool.ErrWrapper{"wrap msg", errors.New("error")}).Error(), "wrap msg: error")
}
