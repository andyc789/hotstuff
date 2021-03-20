// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/relab/hotstuff (interfaces: Config)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	hotstuff "github.com/relab/hotstuff"
)

// MockConfig is a mock of Config interface.
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig.
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance.
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// Fetch mocks base method.
func (m *MockConfig) Fetch(arg0 context.Context, arg1 hotstuff.Hash) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Fetch", arg0, arg1)
}

// Fetch indicates an expected call of Fetch.
func (mr *MockConfigMockRecorder) Fetch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockConfig)(nil).Fetch), arg0, arg1)
}

// Len mocks base method.
func (m *MockConfig) Len() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Len")
	ret0, _ := ret[0].(int)
	return ret0
}

// Len indicates an expected call of Len.
func (mr *MockConfigMockRecorder) Len() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Len", reflect.TypeOf((*MockConfig)(nil).Len))
}

// Propose mocks base method.
func (m *MockConfig) Propose(arg0 *hotstuff.Block) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Propose", arg0)
}

// Propose indicates an expected call of Propose.
func (mr *MockConfigMockRecorder) Propose(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Propose", reflect.TypeOf((*MockConfig)(nil).Propose), arg0)
}

// QuorumSize mocks base method.
func (m *MockConfig) QuorumSize() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QuorumSize")
	ret0, _ := ret[0].(int)
	return ret0
}

// QuorumSize indicates an expected call of QuorumSize.
func (mr *MockConfigMockRecorder) QuorumSize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QuorumSize", reflect.TypeOf((*MockConfig)(nil).QuorumSize))
}

// Replica mocks base method.
func (m *MockConfig) Replica(arg0 hotstuff.ID) (hotstuff.Replica, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Replica", arg0)
	ret0, _ := ret[0].(hotstuff.Replica)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Replica indicates an expected call of Replica.
func (mr *MockConfigMockRecorder) Replica(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Replica", reflect.TypeOf((*MockConfig)(nil).Replica), arg0)
}

// Replicas mocks base method.
func (m *MockConfig) Replicas() map[hotstuff.ID]hotstuff.Replica {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Replicas")
	ret0, _ := ret[0].(map[hotstuff.ID]hotstuff.Replica)
	return ret0
}

// Replicas indicates an expected call of Replicas.
func (mr *MockConfigMockRecorder) Replicas() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Replicas", reflect.TypeOf((*MockConfig)(nil).Replicas))
}

// Timeout mocks base method.
func (m *MockConfig) Timeout(arg0 hotstuff.TimeoutMsg) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Timeout", arg0)
}

// Timeout indicates an expected call of Timeout.
func (mr *MockConfigMockRecorder) Timeout(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Timeout", reflect.TypeOf((*MockConfig)(nil).Timeout), arg0)
}
