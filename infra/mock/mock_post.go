// Code generated by MockGen. DO NOT EDIT.
// Source: post.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

// MockPost is a mock of Post interface.
type MockPost struct {
	ctrl     *gomock.Controller
	recorder *MockPostMockRecorder
}

// MockPostMockRecorder is the mock recorder for MockPost.
type MockPostMockRecorder struct {
	mock *MockPost
}

// NewMockPost creates a new mock instance.
func NewMockPost(ctrl *gomock.Controller) *MockPost {
	mock := &MockPost{ctrl: ctrl}
	mock.recorder = &MockPostMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPost) EXPECT() *MockPostMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockPost) Delete(ctx context.Context, post *entity.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockPostMockRecorder) Delete(ctx, post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPost)(nil).Delete), ctx, post)
}

// FindByID mocks base method.
func (m *MockPost) FindByID(ctx context.Context, postID int) (*entity.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, postID)
	ret0, _ := ret[0].(*entity.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockPostMockRecorder) FindByID(ctx, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockPost)(nil).FindByID), ctx, postID)
}

// FindByUserID mocks base method.
func (m *MockPost) FindByUserID(ctx context.Context, uid string) ([]*entity.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserID", ctx, uid)
	ret0, _ := ret[0].([]*entity.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserID indicates an expected call of FindByUserID.
func (mr *MockPostMockRecorder) FindByUserID(ctx, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserID", reflect.TypeOf((*MockPost)(nil).FindByUserID), ctx, uid)
}

// GetAll mocks base method.
func (m *MockPost) GetAll(ctx context.Context) ([]*entity.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]*entity.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockPostMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockPost)(nil).GetAll), ctx)
}

// Insert mocks base method.
func (m *MockPost) Insert(ctx context.Context, post *entity.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockPostMockRecorder) Insert(ctx, post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockPost)(nil).Insert), ctx, post)
}

// Update mocks base method.
func (m *MockPost) Update(ctx context.Context, post *entity.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockPostMockRecorder) Update(ctx, post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPost)(nil).Update), ctx, post)
}
