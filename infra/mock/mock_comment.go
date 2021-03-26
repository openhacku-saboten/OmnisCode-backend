// Code generated by MockGen. DO NOT EDIT.
// Source: comment.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

// MockComment is a mock of Comment interface.
type MockComment struct {
	ctrl     *gomock.Controller
	recorder *MockCommentMockRecorder
}

// MockCommentMockRecorder is the mock recorder for MockComment.
type MockCommentMockRecorder struct {
	mock *MockComment
}

// NewMockComment creates a new mock instance.
func NewMockComment(ctrl *gomock.Controller) *MockComment {
	mock := &MockComment{ctrl: ctrl}
	mock.recorder = &MockCommentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockComment) EXPECT() *MockCommentMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockComment) Delete(ctx context.Context, userID string, postID, commentID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, userID, postID, commentID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCommentMockRecorder) Delete(ctx, userID, postID, commentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockComment)(nil).Delete), ctx, userID, postID, commentID)
}

// FindByID mocks base method.
func (m *MockComment) FindByID(ctx context.Context, postID, commentID int) (*entity.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, postID, commentID)
	ret0, _ := ret[0].(*entity.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockCommentMockRecorder) FindByID(ctx, postID, commentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockComment)(nil).FindByID), ctx, postID, commentID)
}

// FindByPostID mocks base method.
func (m *MockComment) FindByPostID(ctx context.Context, postID int) ([]*entity.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByPostID", ctx, postID)
	ret0, _ := ret[0].([]*entity.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPostID indicates an expected call of FindByPostID.
func (mr *MockCommentMockRecorder) FindByPostID(ctx, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPostID", reflect.TypeOf((*MockComment)(nil).FindByPostID), ctx, postID)
}

// FindByUserID mocks base method.
func (m *MockComment) FindByUserID(ctx context.Context, uid string) ([]*entity.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserID", ctx, uid)
	ret0, _ := ret[0].([]*entity.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserID indicates an expected call of FindByUserID.
func (mr *MockCommentMockRecorder) FindByUserID(ctx, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserID", reflect.TypeOf((*MockComment)(nil).FindByUserID), ctx, uid)
}

// Insert mocks base method.
func (m *MockComment) Insert(ctx context.Context, comment *entity.Comment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, comment)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockCommentMockRecorder) Insert(ctx, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockComment)(nil).Insert), ctx, comment)
}

// Update mocks base method.
func (m *MockComment) Update(ctx context.Context, comment *entity.Comment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, comment)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCommentMockRecorder) Update(ctx, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockComment)(nil).Update), ctx, comment)
}
