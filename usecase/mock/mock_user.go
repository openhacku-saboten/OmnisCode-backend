// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	entity "github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	reflect "reflect"
)

// MockUser is a mock of User interface
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockUser) EXPECT() *MockUserMockRecorder {
	return _m.recorder
}

// FindByID mocks base method
func (_m *MockUser) FindByID(ctx context.Context, uid string) (*entity.User, error) {
	ret := _m.ctrl.Call(_m, "FindByID", ctx, uid)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID
func (_mr *MockUserMockRecorder) FindByID(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "FindByID", reflect.TypeOf((*MockUser)(nil).FindByID), arg0, arg1)
}

// Insert mocks base method
func (_m *MockUser) Insert(user *entity.User) error {
	ret := _m.ctrl.Call(_m, "Insert", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert
func (_mr *MockUserMockRecorder) Insert(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Insert", reflect.TypeOf((*MockUser)(nil).Insert), arg0)
}

// Update mocks base method
func (_m *MockUser) Update(user *entity.User) error {
	ret := _m.ctrl.Call(_m, "Update", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (_mr *MockUserMockRecorder) Update(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Update", reflect.TypeOf((*MockUser)(nil).Update), arg0)
}
