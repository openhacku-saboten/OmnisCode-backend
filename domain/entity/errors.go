package entity

import (
	"errors"
)

var (
	// ErrUserNotFound はユーザーがDBに存在しないときのエラー
	ErrUserNotFound = errors.New("user not found")

	// ErrDuplicatedUser はユーザーが既にDBに存在するときのエラー
	ErrDuplicatedUser = errors.New("user already exists")

	// ErrDuplicatedUser は入力したTwitterIDが既に使われているときのエラー
	ErrDuplicatedTwitterID = errors.New("twitter id is already used")
)
