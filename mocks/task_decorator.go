// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	pipeline "github.com/egnd/go-pipeline"
	mock "github.com/stretchr/testify/mock"
)

// TaskDecorator is an autogenerated mock type for the TaskDecorator type
type TaskDecorator struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *TaskDecorator) Execute(_a0 pipeline.TaskExecutor) pipeline.TaskExecutor {
	ret := _m.Called(_a0)

	var r0 pipeline.TaskExecutor
	if rf, ok := ret.Get(0).(func(pipeline.TaskExecutor) pipeline.TaskExecutor); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pipeline.TaskExecutor)
		}
	}

	return r0
}

type mockConstructorTestingTNewTaskDecorator interface {
	mock.TestingT
	Cleanup(func())
}

// NewTaskDecorator creates a new instance of TaskDecorator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTaskDecorator(t mockConstructorTestingTNewTaskDecorator) *TaskDecorator {
	mock := &TaskDecorator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
