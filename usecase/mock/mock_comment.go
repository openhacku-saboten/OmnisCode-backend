// Code generated by MockGen. DO NOT EDIT.
// Source: comment.go

// Package mock is a generated GoMock package.
package mock

import (
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

// FindByID mocks base method.
func (m *MockComment) FindByID(postid, commentid int) (*entity.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", postid, commentid)
	ret0, _ := ret[0].(*entity.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockCommentMockRecorder) FindByID(postid, commentid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockComment)(nil).FindByID), postid, commentid)
}

// FindByPostID mocks base method.
func (m *MockComment) FindByPostID(postid int) ([]*entity.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByPostID", postid)
	ret0, _ := ret[0].([]*entity.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPostID indicates an expected call of FindByPostID.
func (mr *MockCommentMockRecorder) FindByPostID(postid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPostID", reflect.TypeOf((*MockComment)(nil).FindByPostID), postid)
}

// Insert mocks base method.
func (m *MockComment) Insert(comment *entity.Comment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", comment)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockCommentMockRecorder) Insert(comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockComment)(nil).Insert), comment)
}
