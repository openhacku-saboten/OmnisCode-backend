package entity_test

import (
	"errors"
	"testing"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

func TestIsValid(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		entity  *entity.Post
		wantErr error
	}{
		{
			name: "正常なPost",
			entity: &entity.Post{
				ID:       "1234",
				UserID:   "1234",
				Title:    "test",
				Code:     "print(\"hello, world\")",
				Language: "Python",
				Content:  "改行しないようにしたい",
				Source:   "",
			},
			wantErr: nil,
		},
		{
			name: "IDが空であるPost",
			entity: &entity.Post{
				ID:       "",
				UserID:   "1234",
				Title:    "test",
				Code:     "print(\"hello, world\")",
				Language: "Python",
				Content:  "改行しないようにしたい",
				Source:   "",
			},
			wantErr: errors.New("post ID must not be empty"),
		},
		{
			name: "userIDが空であるPost",
			entity: &entity.Post{
				ID:       "1234",
				UserID:   "",
				Title:    "test",
				Code:     "print(\"hello, world\")",
				Language: "Python",
				Content:  "改行しないようにしたい",
				Source:   "",
			},
			wantErr: errors.New("post userID must not be empty"),
		},
		{
			name: "Titleが空であるPost",
			entity: &entity.Post{
				ID:       "1234",
				UserID:   "1234",
				Title:    "",
				Code:     "print(\"hello, world\")",
				Language: "Python",
				Content:  "改行しないようにしたい",
				Source:   "",
			},
			wantErr: errors.New("post title must not be empty"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.entity.IsValid()
			if got == nil && tc.wantErr == nil {
				return
			}
			if got.Error() != tc.wantErr.Error() {
				t.Errorf("post.IsValid() = %s, want = %s", got.Error(), tc.wantErr.Error())
			}
		})
	}
}
