package entity

import "errors"

// Post は投稿に関するエンティティです
type Post struct {
	ID       string
	UserID   string
	Title    string
	Code     string
	Language string
	Content  string
	Source   string
}

// IsValid は各エンティティに問題がある場合はerrorを返すメソッドです
func (p *Post) IsValid() error {
	if len(p.ID) == 0 {
		return errors.New("post ID must not be empty")
	}
	if len(p.UserID) == 0 {
		return errors.New("post userID must not be empty")
	}
	if len(p.Title) == 0 {
		return errors.New("post title must not be empty")
	}
	return nil
}
