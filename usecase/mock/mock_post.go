// Code generated by MockGen. DO NOT EDIT.
// Source: post.go

package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	entity "github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	reflect "reflect"
)

// MockPost is a mock of Post interface
type MockPost struct {
	ctrl     *gomock.Controller
	recorder *MockPostMockRecorder
}

// MockPostMockRecorder is the mock recorder for MockPost
type MockPostMockRecorder struct {
	mock *MockPost
}

// NewMockPost creates a new mock instance
func NewMockPost(ctrl *gomock.Controller) *MockPost {
	mock := &MockPost{ctrl: ctrl}
	mock.recorder = &MockPostMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockPost) EXPECT() *MockPostMockRecorder {
	return _m.recorder
}

// GetAll mocks base method
func (_m *MockPost) GetAll(ctx context.Context) ([]*entity.Post, error) {
	ret := _m.ctrl.Call(_m, "GetAll", ctx)
	ret0, _ := ret[0].([]*entity.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (_mr *MockPostMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetAll", reflect.TypeOf((*MockPost)(nil).GetAll), arg0)
}

// Insert mocks base method
func (_m *MockPost) Insert(ctx context.Context, post *entity.Post) error {
	ret := _m.ctrl.Call(_m, "Insert", ctx, post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert
func (_mr *MockPostMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Insert", reflect.TypeOf((*MockPost)(nil).Insert), arg0, arg1)
}