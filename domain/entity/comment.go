package entity

type Comment struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	PostID    int    `json:"post_id"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	FirstLine int    `json:"first_line"`
	LastLine  int    `json:"last_line"`
	Code      string `json:"code"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (c *Comment) IsValid() error {
	if c.ID < 0 {
		return NewErrorNegativeValue("comment ID")
	}
	if len(c.UserID) == 0 {
		return NewErrorEmpty("comment UserID")
	}
	if len([]rune(c.UserID)) > 128 {
		return NewErrorTooLong("comment UserID")
	}
	if c.PostID == 0 {
		return NewErrorEmpty("comment PostID")
	}
	// Typeに応じて必要なフィールドが含まれていなかったらエラー
	switch c.Type {
	case "none":
		// Contentが空ならエラー
		if len(c.Content) == 0 {
			return NewErrorEmpty("comment Content")
		}
	case "highlight":
		// Contentは空でも良い
		// FirstLine,LastLineが空ならエラー
		if c.FirstLine <= 0 {
			return NewErrorEmpty("comment FirstLine")
		}
		if c.LastLine <= 0 {
			return NewErrorEmpty("comment LastLine")
		}
	case "commit":
		// Contentは空でも良い
		// Codeが空ならエラー
		if len(c.Code) == 0 {
			return NewErrorEmpty("comment Code")
		}
	default:
		// none,highlight,commit以外の文字列の場合
		return ErrInvalidCommentType
	}
	return nil
}
