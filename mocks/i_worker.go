// Code generated by mockery v2.10.4. DO NOT EDIT.

package mocks

import (
	wpool "github.com/egnd/go-wpool"
	mock "github.com/stretchr/testify/mock"
)

// IWorker is an autogenerated mock type for the IWorker type
type IWorker struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *IWorker) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Do provides a mock function with given fields: _a0
func (_m *IWorker) Do(_a0 wpool.ITask) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(wpool.ITask) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
