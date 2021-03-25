package entity

import (
	"strings"
	"testing"
)

func TestComment_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		comment *Comment
		wantErr error
	}{
		{
			name: "noneに問題なければnilを返す",
			comment: &Comment{
				ID:      1,
				UserID:  "user-id",
				PostID:  1,
				Type:    "none",
				Content: "type none",
			},
			wantErr: nil,
		},
		{
			name: "highlightに問題なければnilを返す",
			comment: &Comment{
				ID:        1,
				UserID:    "user-id",
				PostID:    1,
				Type:      "highlight",
				Content:   "type highlight",
				FirstLine: 10,
				LastLine:  12,
			},
			wantErr: nil,
		},
		{
			name: "commitに問題なければnilを返す",
			comment: &Comment{
				ID:      1,
				UserID:  "user-id",
				PostID:  1,
				Type:    "commit",
				Content: "type commit",
				Code:    "aaa",
			},
			wantErr: nil,
		},
		{
			name: "IDがマイナスならエラー",
			comment: &Comment{
				ID:      -1,
				UserID:  "user-id",
				PostID:  1,
				Type:    "none",
				Content: "type none",
			},
			wantErr: NewErrorNegativeValue("comment ID"),
		},
		{
			name: "UserIDが空ならエラー",
			comment: &Comment{
				ID:      1,
				PostID:  1,
				Type:    "none",
				Content: "type none",
			},
			wantErr: NewErrorEmpty("comment UserID"),
		},
		{
			name: "UserIDが129文字以上ならエラー",
			comment: &Comment{
				ID:      1,
				UserID:  strings.Repeat("a", 129),
				PostID:  1,
				Type:    "none",
				Content: "type none",
			},
			wantErr: NewErrorTooLong("comment UserID"),
		},
		{
			name: "PostIDが空ならエラー",
			comment: &Comment{
				ID:      1,
				UserID:  "user-id",
				Type:    "none",
				Content: "type none",
			},
			wantErr: NewErrorEmpty("comment PostID"),
		},
		{
			name: "TypeがnoneなのにContentが空ならエラー",
			comment: &Comment{
				ID:     1,
				UserID: "user-id",
				PostID: 1,
				Type:   "none",
			},
			wantErr: NewErrorEmpty("comment Content"),
		},
		{
			name: "TypeがhighlightなのにFirstLineが空ならエラー",
			comment: &Comment{
				ID:      1,
				UserID:  "user-id",
				PostID:  1,
				Type:    "highlight",
				Content: "type highlight",
				Code:    "aaa",
			},
			wantErr: NewErrorEmpty("comment FirstLine"),
		},
		{
			name: "TypeがcommitなのにCodeが空ならエラー",
			comment: &Comment{
				ID:        1,
				UserID:    "user-id",
				PostID:    1,
				Type:      "commit",
				Content:   "type commit",
				FirstLine: 10,
			},
			wantErr: NewErrorEmpty("comment Code"),
		},
		{
			name: "Typeがnone,highlight,commitでなかったらエラー",
			comment: &Comment{
				ID:      1,
				UserID:  "user-id",
				PostID:  1,
				Type:    "invalid",
				Content: "type invalid",
			},
			wantErr: ErrInvalidCommentType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.comment.IsValid()
			if got == nil && tt.wantErr == nil {
				return
			}
			if got.Error() != tt.wantErr.Error() {
				t.Errorf("User.IsValid() = %s, want = %s", got.Error(), tt.wantErr.Error())
			}
		})
	}
}
