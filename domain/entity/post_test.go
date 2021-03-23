package entity_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

func TestIsValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		postE   *entity.Post
		wantErr error
	}{
		{
			name: "正常なentity",
			postE: &entity.Post{
				ID:        0,
				UserID:    "testID",
				Title:     "test title",
				Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language:  "Go",
				Content:   "Test code",
				Source:    "github.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: nil,
		},
		{
			name: "negative value for ID",
			postE: &entity.Post{
				ID:        -1,
				UserID:    "", // empty
				Title:     "test title",
				Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language:  "Go",
				Content:   "Test code",
				Source:    "github.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: errors.New("ID must not be a negative value"),
		},
		{
			name: "empty userID",
			postE: &entity.Post{
				ID:        0,
				UserID:    "", // empty
				Title:     "test title",
				Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language:  "Go",
				Content:   "Test code",
				Source:    "github.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: entity.NewErrorEmpty("post ID"),
		},
		{
			name: "too long userID",
			postE: &entity.Post{
				ID:        0,
				UserID:    strings.Repeat("a", 129), // too long
				Title:     "test title",
				Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language:  "Go",
				Content:   "Test code",
				Source:    "github.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: entity.ErrTooLong,
		},
		{
			name: "empty title",
			postE: &entity.Post{
				ID:        0,
				UserID:    "test",
				Title:     "", // empty
				Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language:  "Go",
				Content:   "Test code",
				Source:    "github.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: entity.NewErrorEmpty("post title"),
		},
		{
			name: "too long title",
			postE: &entity.Post{
				ID:        0,
				UserID:    "test",
				Title:     strings.Repeat("a", 129), // too long
				Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language:  "Go",
				Content:   "Test code",
				Source:    "github.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: entity.ErrTooLong,
		},
		{
			name: "empty code",
			postE: &entity.Post{
				ID:        0,
				UserID:    "test",
				Title:     "test title",
				Code:      "", // empty
				Language:  "Go",
				Content:   "Test code",
				Source:    "github.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: entity.NewErrorEmpty("post code"),
		},
		{
			name: "empty language",
			postE: &entity.Post{
				ID:        0,
				UserID:    "test",
				Title:     "test title",
				Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language:  "", // empty
				Content:   "Test code",
				Source:    "github.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: entity.NewErrorEmpty("post language"),
		},
		{
			name: "too long language",
			postE: &entity.Post{
				ID:        0,
				UserID:    "test",
				Title:     "test title",
				Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language:  strings.Repeat("a", 129), // too long
				Content:   "Test code",
				Source:    "github.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: entity.ErrTooLong,
		},
		{
			name: "too long source",
			postE: &entity.Post{
				ID:        0,
				UserID:    "test",
				Title:     "test title",
				Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language:  "Go",
				Content:   "Test code",
				Source:    strings.Repeat("a", 2049), // too long,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: entity.ErrTooLong,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.postE.IsValid()
			if got == nil && tc.wantErr == nil {
				return
			}
			if got.Error() != tc.wantErr.Error() {
				t.Errorf("postE.IsValid() = %s, want = %s", got.Error(), tc.wantErr.Error())
			}
		})
	}
}
