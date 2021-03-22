package entity

import (
	"errors"
)

var (
	// ErrUserNotFound はユーザーがDBに存在しないときのエラー
	ErrUserNotFound = errors.New("user not found")
	// ErrDuplicatedUser はユーザーが既にDBに存在するときのエラー
	ErrDuplicatedUser = errors.New("user already exists")
	// ErrDuplicatedTwitterID は入力したTwitterIDが既に使われているときのエラー
	ErrDuplicatedTwitterID = errors.New("twitter id is already used")
	// ErrEmptyUserName はユーザー名が空だったときのエラー
	ErrEmptyUserName = errors.New("user name must not be empty")

	// ErrTooLong はフィールドの内容が長すぎるときのエラー
	ErrTooLong = errors.New("too long")
)

// ErrEmptyField はフィールド名が空の時のエラー
type ErrEmpty struct {
	error
	FieldName string
}

func NewErrorEmpty(fieldName string) error {
	return ErrEmpty{
		error:     errors.New("empty"),
		FieldName: fieldName,
	}
}
