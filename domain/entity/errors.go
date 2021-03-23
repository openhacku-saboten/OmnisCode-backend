package entity

import (
	"errors"
	"fmt"
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
)

// ErrTooLong はフィールドの内容が長すぎるときのエラー
type ErrTooLong struct {
	fieldName string
}

// ErrEmptyField はフィールド名が空のときのエラー
type ErrEmpty struct {
	fieldName string
}

// ErrNotFound はエンティティが存在しないときのエラー
type ErrNotFound struct {
	entityName string
}

// ErrDuplicated はエンティティが重複したときのエラー
type ErrDuplicated struct {
	entityName string
}

// NewErrorTooLong はフィールド名が空のときのエラーを生成します
func NewErrorTooLong(fieldName string) error {
	return ErrTooLong{
		fieldName: fieldName,
	}
}

func (e ErrTooLong) Error() string {
	return fmt.Sprintf("%s is too long", e.fieldName)
}

// NewErrorEmpty はフィールド名が空のときのエラーを生成します
func NewErrorEmpty(fieldName string) error {
	return ErrEmpty{
		fieldName: fieldName,
	}
}

func (e ErrEmpty) Error() string {
	return fmt.Sprintf("%s is empty", e.fieldName)
}

// NewErrorNotFound はフィールド名が存在しないときのエラーを生成します
func NewErrorNotFound(entityName string) error {
	return ErrNotFound{
		entityName: entityName,
	}
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%s is not found", e.entityName)
}

// NewDuplicated はフィールド名が重複したときのエラーを生成します
func NewErrorDuplicated(entityName string) error {
	return ErrDuplicated{
		entityName: entityName,
	}
}

func (e ErrDuplicated) Error() string {
	return fmt.Sprintf("%s is duplicated", e.entityName)
}
