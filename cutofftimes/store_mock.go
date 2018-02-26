package cutofftimes

// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mock_cutofftimes is a generated GoMock package.

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// saveCutoff mocks base method
func (m *MockStore) saveCutoff(time cutoffTime) {
	m.ctrl.Call(m, "saveCutoff", time)
}

// saveCutoff indicates an expected call of saveCutoff
func (mr *MockStoreMockRecorder) saveCutoff(time interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "saveCutoff", reflect.TypeOf((*MockStore)(nil).saveCutoff), time)
}

// cutoffExists mocks base method
func (m *MockStore) cutoffExists(service, subservice string) bool {
	ret := m.ctrl.Call(m, "cutoffExists", service, subservice)
	ret0, _ := ret[0].(bool)
	return ret0
}

// cutoffExists indicates an expected call of cutoffExists
func (mr *MockStoreMockRecorder) cutoffExists(service, subservice interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "cutoffExists", reflect.TypeOf((*MockStore)(nil).cutoffExists), service, subservice)
}

// isInStartOfDay mocks base method
func (m *MockStore) isInStartOfDay(service, subservice string) bool {
	ret := m.ctrl.Call(m, "isInStartOfDay", service, subservice)
	ret0, _ := ret[0].(bool)
	return ret0
}

// isInStartOfDay indicates an expected call of isInStartOfDay
func (mr *MockStoreMockRecorder) isInStartOfDay(service, subservice interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isInStartOfDay", reflect.TypeOf((*MockStore)(nil).isInStartOfDay), service, subservice)
}
