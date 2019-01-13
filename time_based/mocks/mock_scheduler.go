// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/abdulrahmank/job_executor/time_based/scheduler (interfaces: TimeBasedScheduler)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockTimeBasedScheduler is a mock of TimeBasedScheduler interface
type MockTimeBasedScheduler struct {
	ctrl     *gomock.Controller
	recorder *MockTimeBasedSchedulerMockRecorder
}

// MockTimeBasedSchedulerMockRecorder is the mock recorder for MockTimeBasedScheduler
type MockTimeBasedSchedulerMockRecorder struct {
	mock *MockTimeBasedScheduler
}

// NewMockTimeBasedScheduler creates a new mock instance
func NewMockTimeBasedScheduler(ctrl *gomock.Controller) *MockTimeBasedScheduler {
	mock := &MockTimeBasedScheduler{ctrl: ctrl}
	mock.recorder = &MockTimeBasedSchedulerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTimeBasedScheduler) EXPECT() *MockTimeBasedSchedulerMockRecorder {
	return m.recorder
}

// Schedule mocks base method
func (m *MockTimeBasedScheduler) Schedule(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Schedule", arg0, arg1)
}

// Schedule indicates an expected call of Schedule
func (mr *MockTimeBasedSchedulerMockRecorder) Schedule(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Schedule", reflect.TypeOf((*MockTimeBasedScheduler)(nil).Schedule), arg0, arg1)
}
