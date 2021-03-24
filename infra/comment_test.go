package infra

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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

	dbMap.AddTableWithName(PostDTO{}, "posts")
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
			Type:      "code",
			Content:   "type code",
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
					CreatedAt: "1970-01-01T09:01:40+09:00",
					UpdatedAt: "1970-01-01T09:01:40+09:00",
				},
				{
					ID:        2,
					UserID:    "user-id",
					PostID:    1,
					Type:      "highlight",
					Content:   "type highlight",
					FirstLine: 10,
					LastLine:  11,
					CreatedAt: "1970-01-01T09:01:40+09:00",
					UpdatedAt: "1970-01-01T09:01:40+09:00",
				},
			},
			wantErr: nil,
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
