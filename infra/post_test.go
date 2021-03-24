package infra

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/service"
)

func TestPostRepository_FindByID(t *testing.T) {
	dbMap, err := NewDB()
	if err != nil {
		t.Fatalf(err.Error())
	}
	dbMap.AddTableWithName(PostDTO{}, "posts")
	truncateTable(t, dbMap, "posts")

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

	validPost := &PostDTO{
		ID:        1,
		UserID:    "user-id",
		Title:     "test title",
		Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
		Language:  "Go",
		Content:   "Test code",
		Source:    "github.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// デフォルトの投稿追加
	if err := dbMap.Insert(validPost); err != nil {
		t.Fatal(err)
	}

	postRepo := NewPostRepository(dbMap)

	tests := []struct {
		name     string
		postID   int
		wantPost *entity.Post
		wantErr  error
	}{
		{
			name:   "正常に取得できる",
			postID: 1,
			wantPost: &entity.Post{
				ID:        validPost.ID,
				UserID:    validPost.UserID,
				Title:     validPost.Title,
				Code:      validPost.Code,
				Language:  validPost.Language,
				Content:   validPost.Content,
				Source:    validPost.Source,
				CreatedAt: service.ConvertTimeToStr(validPost.CreatedAt),
				UpdatedAt: service.ConvertTimeToStr(validPost.UpdatedAt),
			},
			wantErr: nil,
		},
		{
			name:     "存在しないユーザは取得できない",
			postID:   0,
			wantPost: nil,
			wantErr:  entity.NewErrorNotFound("post"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := postRepo.FindByID(ctx, tt.postID)

			if err == nil || tt.wantErr == nil {
				if err == tt.wantErr {
					return
				}
				// どちらかがnilの場合は%vを使う
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
			} else if err.Error() != tt.wantErr.Error() {
				t.Errorf("error = %s, wantErr = %s", err.Error(), tt.wantErr.Error())
			}

			if diff := cmp.Diff(got, tt.wantPost); diff != "" {
				t.Errorf("post FindByID(%d) = %+v, want = %+v", tt.postID, got, tt.wantPost)
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
		CreatedAt: time.Now(), // とりあえず入れる
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
				ID:       2,
				UserID:   "user-id",
				Title:    "test title",
				Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language: "Go",
				Content:  "Test code",
				Source:   "github.com",
			},
			wantErr: nil,
		},
		{
			name: "存在しないユーザで登録するとエラー",
			post: &entity.Post{
				ID:       3,
				UserID:   "user-id2",
				Title:    "test title",
				Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language: "Go",
				Content:  "Test code",
				Source:   "github.com",
			},
			wantErr: errors.New("unexisted user"),
		},
		{
			name: "重複したpostIDで登録するとエラー",
			post: &entity.Post{
				ID:       1,
				UserID:   "user-id",
				Title:    "test title",
				Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language: "Go",
				Content:  "Test code",
				Source:   "github.com",
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
