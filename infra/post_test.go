package infra

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

func TestPostRepository_Insert(t *testing.T) {
	dbMap, err := NewDB()
	if err != nil {
		t.Fatalf(err.Error())
	}

	dbMap.AddTableWithName(UserDTO{}, "users")
	truncateTable(t, dbMap, "users")

	if err := dbMap.Insert(&UserDTO{
		ID:        "user-id",
		Name:      "test user",
		Profile:   "test profile",
		TwitterID: "twitter",
	}); err != nil {
		t.Fatal(err)
	}

	dbMap.AddTableWithName(PostDTO{}, "posts")
	truncateTable(t, dbMap, "posts")

	// デフォルトユーザの追加
	if err := dbMap.Insert(&PostDTO{
		ID:        1,
		UserID:    "user-id",
		Title:     "test title",
		Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
		Language:  "Go",
		Content:   "Test code",
		Source:    "github.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		t.Fatal(err)
	}

	postRepo := NewPostRepository(dbMap)

	tests := []struct {
		name    string
		post    *entity.Post
		wantErr error
	}{
		{
			name: "正常に追加できる",
			post: &entity.Post{
				ID:        2,
				UserID:    "user-id",
				Title:     "test title",
				Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language:  "Go",
				Content:   "Test code",
				Source:    "github.com",
				CreatedAt: "2021-03-23T11:42:56+09:00",
				UpdatedAt: "2021-03-23T11:42:56+09:00",
			},
			wantErr: nil,
		},
		{
			name: "存在しないユーザで登録するとエラー",
			post: &entity.Post{
				ID:        3,
				UserID:    "user-id2",
				Title:     "test title",
				Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language:  "Go",
				Content:   "Test code",
				Source:    "github.com",
				CreatedAt: "2021-03-23T11:42:56+09:00",
				UpdatedAt: "2021-03-23T11:42:56+09:00",
			},
			wantErr: errors.New("unexisted user"),
		},
		{
			name: "重複したpostIDで登録するとエラー",
			post: &entity.Post{
				ID:        1,
				UserID:    "user-id",
				Title:     "test title",
				Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language:  "Go",
				Content:   "Test code",
				Source:    "github.com",
				CreatedAt: "2021-03-23T11:42:56+09:00",
				UpdatedAt: "2021-03-23T11:42:56+09:00",
			},
			wantErr: errors.New("post ID is duplicated"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := postRepo.Insert(ctx, tt.post)

			if err == nil || tt.wantErr == nil {
				if err == tt.wantErr {
					return
				}
				// どちらかがnilの場合は%vを使う
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
			} else if err.Error() != tt.wantErr.Error() {
				t.Errorf("error = %s, wantErr = %s", err.Error(), tt.wantErr.Error())
			}
		})
	}
}
