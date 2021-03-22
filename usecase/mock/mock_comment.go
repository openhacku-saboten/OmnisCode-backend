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

// GetByPostID mocks base method.
func (m *MockComment) GetByPostID(postid int) ([]*entity.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByPostID", postid)
	ret0, _ := ret[0].([]*entity.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByPostID indicates an expected call of GetByPostID.
func (mr *MockCommentMockRecorder) GetByPostID(postid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByPostID", reflect.TypeOf((*MockComment)(nil).GetByPostID), postid)
}
