package infra

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

func TestCommentRepository_GetByPostID(t *testing.T) {
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

	dbMap.AddTableWithName(PostDTO{}, "posts").SetKeys(true, "id")
	truncateTable(t, dbMap, "posts")

	postDTOs := []*PostDTO{
		{
			ID:        1,
			UserID:    "user-id",
			Title:     "test title",
			Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
			Language:  "Go",
			Content:   "Test code",
			Source:    "github.com",
			CreatedAt: time.Unix(100, 0), // とりあえず入れる
			UpdatedAt: time.Unix(100, 0),
		},
		{
			ID:        2,
			UserID:    "user-id",
			Title:     "test title",
			Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
			Language:  "Go",
			Content:   "Test code",
			Source:    "github.com",
			CreatedAt: time.Unix(100, 0), // とりあえず入れる
			UpdatedAt: time.Unix(100, 0),
		},
	}
	for _, postDTO := range postDTOs {
		if err := dbMap.Insert(postDTO); err != nil {
			t.Fatal(err)
		}
	}

	dbMap.AddTableWithName(CommentDTO{}, "comments")
	truncateTable(t, dbMap, "comments")

	commentDTOs := []*CommentDTO{
		{
			ID:        1,
			UserID:    "user-id",
			PostID:    1,
			Type:      "none",
			Content:   "type none",
			CreatedAt: time.Unix(100, 0), // とりあえず入れる
			UpdatedAt: time.Unix(100, 0),
		},
		{
			ID:        2,
			UserID:    "user-id",
			PostID:    1,
			Type:      "highlight",
			Content:   "type highlight",
			FirstLine: 10,
			LastLine:  11,
			CreatedAt: time.Unix(100, 0), // とりあえず入れる
			UpdatedAt: time.Unix(100, 0),
		},
		{
			ID:        3,
			UserID:    "user-id",
			PostID:    2,
			Type:      "commit",
			Content:   "type commit",
			Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
			CreatedAt: time.Unix(100, 0), // とりあえず入れる
			UpdatedAt: time.Unix(100, 0),
		},
	}
	for _, commentDTO := range commentDTOs {
		if err := dbMap.Insert(commentDTO); err != nil {
			t.Fatal(err)
		}
	}

	commentRepo := NewCommentRepository(dbMap)

	tests := []struct {
		name         string
		postID       int
		wantComments []*entity.Comment
		wantErr      error
	}{
		{
			name:   "正しくコメントを取得できる",
			postID: 1,
			wantComments: []*entity.Comment{
				{
					ID:        1,
					UserID:    "user-id",
					PostID:    1,
					Type:      "none",
					Content:   "type none",
					CreatedAt: "1970-01-01T00:01:40+09:00",
					UpdatedAt: "1970-01-01T00:01:40+09:00",
				},
				{
					ID:        2,
					UserID:    "user-id",
					PostID:    1,
					Type:      "highlight",
					Content:   "type highlight",
					FirstLine: 10,
					LastLine:  11,
					CreatedAt: "1970-01-01T00:01:40+09:00",
					UpdatedAt: "1970-01-01T00:01:40+09:00",
				},
			},
			wantErr: nil,
		},
		{
			name:         "コメントが存在しなければErrNotFound",
			postID:       100,
			wantComments: nil,
			wantErr:      entity.NewErrorNotFound("comment"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotComments, err := commentRepo.GetByPostID(tt.postID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				diff := cmp.Diff(tt.wantComments, gotComments)
				if diff != "" {
					t.Errorf("Data (-want +got) =\n%s\n", diff)
				}
			}
		})
	}
}

func TestUserRepository_GetByUserID(t *testing.T) {
	dbMap, err := NewDB()
	if err != nil {
		t.Fatalf(err.Error())
	}

	dbMap.AddTableWithName(UserDTO{}, "users")
	truncateTable(t, dbMap, "users")

	userDTOs := []*UserDTO{
		{
			ID:        "user-id",
			Name:      "test user",
			Profile:   "test profile",
			TwitterID: "twitter",
		},
		{
			ID:        "user-id2",
			Name:      "test user",
			Profile:   "test profile",
			TwitterID: "twitter2",
		},
	}
	for _, userDTO := range userDTOs {
		if err := dbMap.Insert(userDTO); err != nil {
			t.Fatal(err)
		}
	}

	dbMap.AddTableWithName(PostDTO{}, "posts").SetKeys(true, "id")
	dbMap.AddTableWithName(PostInsertDTO{}, "posts").SetKeys(true, "id")
	truncateTable(t, dbMap, "posts")

	postDTOs := []*PostInsertDTO{
		{
			ID:       1,
			UserID:   "user-id",
			Title:    "test title",
			Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
			Language: "Go",
			Content:  "Test code",
			Source:   "github.com",
		},
		{
			ID:       2,
			UserID:   "user-id",
			Title:    "test title",
			Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
			Language: "Go",
			Content:  "Test code",
			Source:   "github.com",
		},
	}
	for _, postDTO := range postDTOs {
		if err := dbMap.Insert(postDTO); err != nil {
			t.Fatal(err)
		}
	}

	dbMap.AddTableWithName(CommentDTO{}, "comments").SetKeys(true, "id")
	dbMap.AddTableWithName(CommentInsertDTO{}, "comments").SetKeys(true, "id")
	truncateTable(t, dbMap, "comments")

	commentDTOs := []*CommentInsertDTO{
		{
			ID:      1,
			UserID:  "user-id",
			PostID:  1,
			Type:    "none",
			Content: "type none",
		},
		{
			ID:        2,
			UserID:    "user-id2",
			PostID:    1,
			Type:      "highlight",
			Content:   "type highlight",
			FirstLine: 10,
			LastLine:  11,
		},
		{
			ID:      3,
			UserID:  "user-id",
			PostID:  2,
			Type:    "commit",
			Content: "type commit",
			Code:    "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
		},
	}
	for _, commentDTO := range commentDTOs {
		if err := dbMap.Insert(commentDTO); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name         string
		userID       string
		wantComments []*entity.Comment
		wantErr      error
	}{
		{
			name:   "正しくコメントを取得できる",
			userID: "user-id",
			wantComments: []*entity.Comment{
				{
					ID:      1,
					UserID:  "user-id",
					PostID:  1,
					Type:    "none",
					Content: "type none",
				},
				{
					ID:      3,
					UserID:  "user-id",
					PostID:  2,
					Type:    "commit",
					Content: "type commit",
					Code:    "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				},
			},
			wantErr: nil,
		},
		{
			name:         "コメントが存在しなければErrNotFound",
			userID:       "user-id100",
			wantComments: nil,
			wantErr:      entity.NewErrorNotFound("comment"),
		},
	}

	commentRepo := NewCommentRepository(dbMap)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			comments, err := commentRepo.FindByUserID(ctx, tt.userID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				for idx := range tt.wantComments {
					diff := cmp.Diff(tt.wantComments[idx], comments[idx], cmpopts.IgnoreFields(entity.Comment{}, "CreatedAt", "UpdatedAt"))
					if diff != "" {
						t.Errorf("Data (-want +got) =\n%s\n", diff)
					}
				}
			}
		})
	}
}
