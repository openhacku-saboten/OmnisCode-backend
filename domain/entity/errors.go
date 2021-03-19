package entity

import "errors"

var (
	// ErrInvalidUser はUserのフィールドに不正があるエラーを示す
	ErrInvalidUser = errors.New("Invalid user")
)
