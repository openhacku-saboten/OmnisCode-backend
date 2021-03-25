package entity

import (
	"errors"
	"strconv"
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

// PostForJSON はJSONにBindする時のPostです
type PostForJSON struct {
	ID        string `json:"id"`
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

// Convert はPostForJSONからPostに変換します
func (p PostForJSON) Convert() (*Post, error) {
	id, err := strconv.Atoi(p.ID)
	if err != nil {
		return nil, err
	}

	return &Post{
		ID:        id,
		UserID:    p.UserID, // c.Get("userID")で更新される
		Title:     p.Title,
		Code:      p.Code,
		Language:  p.Language,
		Content:   p.Content,
		Source:    p.Source,
		CreatedAt: p.CreatedAt, // 永続化する時に自動更新される
		UpdatedAt: p.UpdatedAt, // 永続化する時に自動更新される
	}, nil
}

// Convert はPostからPostForJSONに変換します
func (p Post) Convert() *PostForJSON {
	return &PostForJSON{
		ID:        strconv.Itoa(p.ID),
		UserID:    p.UserID, // c.Get("userID")で更新される
		Title:     p.Title,
		Code:      p.Code,
		Language:  p.Language,
		Content:   p.Content,
		Source:    p.Source,
		CreatedAt: p.CreatedAt, // 永続化する時に自動更新される
		UpdatedAt: p.UpdatedAt, // 永続化する時に自動更新される
	}
}
