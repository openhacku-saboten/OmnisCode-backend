package infra

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

func TestPostRepository_GetAll(t *testing.T) {
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

	post1 := &entity.Post{
		ID:        1,
		UserID:    "user-id",
		Title:     "test title",
		Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
		Language:  "Go",
		Content:   "Test code",
		Source:    "github.com",
		CreatedAt: "2021-03-23T11:42:56+09:00",
		UpdatedAt: "2021-03-23T11:42:56+09:00",
	}

	post2 := &entity.Post{
		ID:        2,
		UserID:    "user-id",
		Title:     "test title",
		Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
		Language:  "Go",
		Content:   "Test code",
		Source:    "github.com",
		CreatedAt: "2021-03-23T11:42:56+09:00",
		UpdatedAt: "2021-03-23T11:42:56+09:00",
	}

	postRepo := NewPostRepository(dbMap)

	tests := []struct {
		name          string
		insertedPosts []*entity.Post
		wantPosts     []*entity.Post
		wantErr       error
	}{
		{
			name:          "正常に取得できる",
			insertedPosts: []*entity.Post{post1, post2},
			wantPosts:     []*entity.Post{post1, post2},
			wantErr:       nil,
		},
		{
			name:          "投稿がない時はNotFound",
			insertedPosts: []*entity.Post{post1, post2},
			wantPosts:     []*entity.Post{post1, post2},
			wantErr:       entity.NewErrorNotFound("post"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			for _, post := range tt.insertedPosts {
				if err := postRepo.Insert(ctx, post); err != nil {
					posts, _ := postRepo.GetAll(ctx)
					t.Errorf("%+v", posts)
					t.Fatalf("failed Insert: %s", err.Error())
				}
			}

			posts, err := postRepo.GetAll(ctx)
			if err == nil || tt.wantErr == nil {
				if err == tt.wantErr {
					return
				}
				// どちらかがnilの場合は%vを使う
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
			} else if err.Error() != tt.wantErr.Error() {
				t.Errorf("error = %s, wantErr = %s", err.Error(), tt.wantErr.Error())
			}

			if diff := cmp.Diff(posts, tt.wantPosts); diff != "" {
				t.Errorf("post GetAll():\n%s", diff)
			}
		})
	}
}

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

	// デフォルトの投稿追加
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
