// Code generated by MockGen. DO NOT EDIT.
// Source: app/cmd/markdown/yaml/configyaml_test.go
//
// Generated by this command:
//
//	mockgen --source app/cmd/markdown/yaml/configyaml_test.go --package yaml
//

// Package yaml is a generated GoMock package.
package files

import (
	reflect "reflect"

	templates "github.com/gSchool/glearn-cli/app/cmd/markdown/templates"
	cobra "github.com/spf13/cobra"
	gomock "go.uber.org/mock/gomock"
)

// MockRunCallback is a mock of RunCallback interface.
type MockRunCallback struct {
	ctrl     *gomock.Controller
	recorder *MockRunCallbackMockRecorder
}

// MockRunCallbackMockRecorder is the mock recorder for MockRunCallback.
type MockRunCallbackMockRecorder struct {
	mock *MockRunCallback
}

// NewMockRunCallback creates a new mock instance.
func NewMockRunCallback(ctrl *gomock.Controller) *MockRunCallback {
	mock := &MockRunCallback{ctrl: ctrl}
	mock.recorder = &MockRunCallbackMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRunCallback) EXPECT() *MockRunCallbackMockRecorder {
	return m.recorder
}

// Call mocks base method.
func (m *MockRunCallback) Call(arg0 *cobra.Command, arg1 *string, arg2 templates.Template) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Call", arg0, arg1, arg2)
}

// Call indicates an expected call of Call.
func (mr *MockRunCallbackMockRecorder) Call(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockRunCallback)(nil).Call), arg0, arg1, arg2)
}

// MockValidator is a mock of Validator interface.
type MockValidator struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorMockRecorder
}

// MockValidatorMockRecorder is the mock recorder for MockValidator.
type MockValidatorMockRecorder struct {
	mock *MockValidator
}

// NewMockValidator creates a new mock instance.
func NewMockValidator(ctrl *gomock.Controller) *MockValidator {
	mock := &MockValidator{ctrl: ctrl}
	mock.recorder = &MockValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockValidator) EXPECT() *MockValidatorMockRecorder {
	return m.recorder
}

// Call mocks base method.
func (m *MockValidator) Call(arg0 *cobra.Command, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Call indicates an expected call of Call.
func (mr *MockValidatorMockRecorder) Call(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockValidator)(nil).Call), arg0, arg1)
}