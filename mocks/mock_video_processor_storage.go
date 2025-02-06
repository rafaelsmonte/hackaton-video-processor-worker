// Code generated by MockGen. DO NOT EDIT.
// Source: hackaton-video-processor-worker/internal/domain/adapters (interfaces: IVideoProcessorStorage)

// Package mocks is a generated GoMock package.
package mocks

import (
	entities "hackaton-video-processor-worker/internal/domain/entities"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIVideoProcessorStorage is a mock of IVideoProcessorStorage interface.
type MockIVideoProcessorStorage struct {
	ctrl     *gomock.Controller
	recorder *MockIVideoProcessorStorageMockRecorder
}

// MockIVideoProcessorStorageMockRecorder is the mock recorder for MockIVideoProcessorStorage.
type MockIVideoProcessorStorageMockRecorder struct {
	mock *MockIVideoProcessorStorage
}

// NewMockIVideoProcessorStorage creates a new mock instance.
func NewMockIVideoProcessorStorage(ctrl *gomock.Controller) *MockIVideoProcessorStorage {
	mock := &MockIVideoProcessorStorage{ctrl: ctrl}
	mock.recorder = &MockIVideoProcessorStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIVideoProcessorStorage) EXPECT() *MockIVideoProcessorStorageMockRecorder {
	return m.recorder
}

// Download mocks base method.
func (m *MockIVideoProcessorStorage) Download(arg0 entities.File) (entities.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", arg0)
	ret0, _ := ret[0].(entities.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download.
func (mr *MockIVideoProcessorStorageMockRecorder) Download(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockIVideoProcessorStorage)(nil).Download), arg0)
}

// Upload mocks base method.
func (m *MockIVideoProcessorStorage) Upload(arg0 entities.File) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upload indicates an expected call of Upload.
func (mr *MockIVideoProcessorStorageMockRecorder) Upload(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*MockIVideoProcessorStorage)(nil).Upload), arg0)
}
