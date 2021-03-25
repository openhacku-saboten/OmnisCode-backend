package entity

import (
	"errors"
)

// Post は投稿を表します
type Post struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Code      string `json:"code"`
	Language  string `json:"language"`
	Content   string `json:"content"`
	Source    string `json:"source"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// IsValid は各エンティティに問題がある場合はerrorを返すメソッドです
func (p *Post) IsValid() error {
	if p.ID < 0 {
		return errors.New("ID must not be a negative value")
	}
	if len(p.UserID) == 0 {
		return NewErrorEmpty("post UserID")
	}
	if len([]rune(p.UserID)) > 128 {
		return NewErrorTooLong("post UserID")
	}
	if len(p.Title) == 0 {
		return NewErrorEmpty("post Title")
	}
	if len([]rune(p.Title)) > 128 {
		return NewErrorTooLong("post Title")
	}
	if len(p.Code) == 0 {
		return NewErrorEmpty("post Code")
	}
	if len(p.Language) == 0 {
		return NewErrorEmpty("post Language")
	}
	if len([]rune(p.Language)) > 128 {
		return NewErrorTooLong("post Language")
	}
	if len([]rune(p.Source)) > 2048 {
		return NewErrorTooLong("post Source")
	}

	return nil

}
