package entity

import (
	"errors"
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	PostID    int       `json:"post_id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	FirstLine int       `json:"first_line"`
	LastLine  int       `json:"last_line"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Comment) IsValid() error {
	if c.ID == 0 {
		return errors.New("comment ID must not be empty")
	}
	if len(c.UserID) == 0 {
		return errors.New("user ID must not be empty")
	}
	if c.PostID == 0 {
		return errors.New("post ID must not be empty")
	}
	switch c.Type {
	case "none":
		if c.FirstLine != 0 || c.LastLine != 0 || len(c.Code) != 0 {
			return errors.New("firstline, lastline and code must be empty if type is none")
		}
	case "highlight":
		if len(c.Code) != 0 {
			return errors.New("code must be empty if type is highlight")
		}
	case "commit":
		if len(c.Code) != 0 {
			return errors.New("firstline and lastline must be empty if type is commit")
		}
	default:
		return errors.New("invalid comment type")
	}
	return nil
}
