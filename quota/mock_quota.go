// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/google/trillian/quota (interfaces: Manager)

// Package quota is a generated GoMock package.
package quota

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// GetTokens mocks base method.
func (m *MockManager) GetTokens(arg0 context.Context, arg1 int, arg2 []Spec) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokens", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetTokens indicates an expected call of GetTokens.
func (mr *MockManagerMockRecorder) GetTokens(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokens", reflect.TypeOf((*MockManager)(nil).GetTokens), arg0, arg1, arg2)
}

// PutTokens mocks base method.
func (m *MockManager) PutTokens(arg0 context.Context, arg1 int, arg2 []Spec) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutTokens", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutTokens indicates an expected call of PutTokens.
func (mr *MockManagerMockRecorder) PutTokens(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutTokens", reflect.TypeOf((*MockManager)(nil).PutTokens), arg0, arg1, arg2)
}

// ResetQuota mocks base method.
func (m *MockManager) ResetQuota(arg0 context.Context, arg1 []Spec) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetQuota", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetQuota indicates an expected call of ResetQuota.
func (mr *MockManagerMockRecorder) ResetQuota(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetQuota", reflect.TypeOf((*MockManager)(nil).ResetQuota), arg0, arg1)
}
