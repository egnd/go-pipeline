package pkg_test

import (
	"errors"
	"testing"

	"github.com/egnd/wpool/pkg"
	"github.com/stretchr/testify/assert"
)

func Test_ErrIsClosed(t *testing.T) {
	assert.EqualValues(t, (&pkg.ErrIsClosed{"pool"}).Error(), "pool is closed")
}

func Test_ErrPanic(t *testing.T) {
	assert.EqualValues(t, (&pkg.ErrPanic{"panic"}).Error(), "panic occurred: panic")
}

func Test_ErrWrapper(t *testing.T) {
	assert.EqualValues(t, (&pkg.ErrWrapper{"wrap msg", errors.New("error")}).Error(), "wrap msg: error")
}
