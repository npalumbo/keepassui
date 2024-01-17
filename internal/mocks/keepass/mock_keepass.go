// Code generated by MockGen. DO NOT EDIT.
// Source: internal/keepass/keepass.go
//
// Generated by this command:
//
//	mockgen -source=internal/keepass/keepass.go -destination=internal/keepass/mock_keepass.go
//

// Package mock_keepass is a generated GoMock package.
package mock_keepass

import (
	keepass "keepassui/internal/keepass"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockSecretReader is a mock of SecretReader interface.
type MockSecretReader struct {
	ctrl     *gomock.Controller
	recorder *MockSecretReaderMockRecorder
}

// MockSecretReaderMockRecorder is the mock recorder for MockSecretReader.
type MockSecretReaderMockRecorder struct {
	mock *MockSecretReader
}

// NewMockSecretReader creates a new mock instance.
func NewMockSecretReader(ctrl *gomock.Controller) *MockSecretReader {
	mock := &MockSecretReader{ctrl: ctrl}
	mock.recorder = &MockSecretReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretReader) EXPECT() *MockSecretReaderMockRecorder {
	return m.recorder
}

// ReadEntriesFromContentGroupedByPath mocks base method.
func (m *MockSecretReader) ReadEntriesFromContentGroupedByPath() (map[string][]keepass.SecretEntry, []string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadEntriesFromContentGroupedByPath")
	ret0, _ := ret[0].(map[string][]keepass.SecretEntry)
	ret1, _ := ret[1].([]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReadEntriesFromContentGroupedByPath indicates an expected call of ReadEntriesFromContentGroupedByPath.
func (mr *MockSecretReaderMockRecorder) ReadEntriesFromContentGroupedByPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadEntriesFromContentGroupedByPath", reflect.TypeOf((*MockSecretReader)(nil).ReadEntriesFromContentGroupedByPath))
}
