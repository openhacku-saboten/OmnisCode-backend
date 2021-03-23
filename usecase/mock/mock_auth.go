// Code generated by MockGen. DO NOT EDIT.
// Source: auth.go

package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAuth is a mock of Auth interface
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return _m.recorder
}

// Authenticate mocks base method
func (_m *MockAuth) Authenticate(ctx context.Context, token string) (string, error) {
	ret := _m.ctrl.Call(_m, "Authenticate", ctx, token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authenticate indicates an expected call of Authenticate
func (_mr *MockAuthMockRecorder) Authenticate(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Authenticate", reflect.TypeOf((*MockAuth)(nil).Authenticate), arg0, arg1)
}

// GetIconURL mocks base method
func (_m *MockAuth) GetIconURL(ctx context.Context, uid string) (string, error) {
	ret := _m.ctrl.Call(_m, "GetIconURL", ctx, uid)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIconURL indicates an expected call of GetIconURL
func (_mr *MockAuthMockRecorder) GetIconURL(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetIconURL", reflect.TypeOf((*MockAuth)(nil).GetIconURL), arg0, arg1)
}
