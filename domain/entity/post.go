package entity

import (
	"errors"
	"time"
)

type Post struct {
	ID        int
	UserID    string
	Title     string
	Code      string
	Language  string
	Content   string
	Source    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsValid は各エンティティに問題がある場合はerrorを返すメソッドです
func (p *Post) IsValid() error {
	if p.ID < 0 {
		return errors.New("ID must not be a negative value")
	}
	if len(p.UserID) == 0 {
		return NewErrorEmpty("post userID")
	}
	if len([]rune(p.UserID)) > 128 {
		return ErrTooLong
	}
	if len(p.Title) == 0 {
		return NewErrorEmpty("post title")
	}
	if len([]rune(p.Title)) > 128 {
		return ErrTooLong
	}
	if len(p.Code) == 0 {
		return NewErrorEmpty("post code")
	}
	if len(p.Language) == 0 {
		return NewErrorEmpty("post language")
	}
	if len([]rune(p.Language)) > 128 {
		return ErrTooLong
	}
	if len([]rune(p.Source)) > 2048 {
		return ErrTooLong
	}

	return nil

}
