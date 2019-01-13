// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/abdulrahmank/job_executor/time_based/dao (interfaces: JobSettingDao)

// Package mocks is a generated GoMock package.
package mocks

import (
	dao "github.com/abdulrahmank/job_executor/time_based/dao"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockJobSettingDao is a mock of JobSettingDao interface
type MockJobSettingDao struct {
	ctrl     *gomock.Controller
	recorder *MockJobSettingDaoMockRecorder
}

// MockJobSettingDaoMockRecorder is the mock recorder for MockJobSettingDao
type MockJobSettingDaoMockRecorder struct {
	mock *MockJobSettingDao
}

// NewMockJobSettingDao creates a new mock instance
func NewMockJobSettingDao(ctrl *gomock.Controller) *MockJobSettingDao {
	mock := &MockJobSettingDao{ctrl: ctrl}
	mock.recorder = &MockJobSettingDaoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockJobSettingDao) EXPECT() *MockJobSettingDaoMockRecorder {
	return m.recorder
}

// GetJobFor mocks base method
func (m *MockJobSettingDao) GetJobFor(arg0 string) []dao.JobSettings {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJobFor", arg0)
	ret0, _ := ret[0].([]dao.JobSettings)
	return ret0
}

// GetJobFor indicates an expected call of GetJobFor
func (mr *MockJobSettingDaoMockRecorder) GetJobFor(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJobFor", reflect.TypeOf((*MockJobSettingDao)(nil).GetJobFor), arg0)
}

// SaveJob mocks base method
func (m *MockJobSettingDao) SaveJob(arg0, arg1, arg2, arg3 string, arg4 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveJob", arg0, arg1, arg2, arg3, arg4)
}

// SaveJob indicates an expected call of SaveJob
func (mr *MockJobSettingDaoMockRecorder) SaveJob(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveJob", reflect.TypeOf((*MockJobSettingDao)(nil).SaveJob), arg0, arg1, arg2, arg3, arg4)
}
