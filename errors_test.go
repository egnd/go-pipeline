package wpool_test

import (
	"errors"
	"testing"

	"github.com/egnd/go-wpool"
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

func Test_ErrTaskTimeout(t *testing.T) {
	assert.EqualValues(t, (&wpool.ErrTaskTimeout{"taskname"}).Error(), "task timeout: taskname")
	assert.EqualValues(t, (&wpool.ErrTaskTimeout{"taskname"}).Timeout(), true)
	assert.EqualValues(t, (&wpool.ErrTaskTimeout{"taskname"}).Temporary(), true)
}
