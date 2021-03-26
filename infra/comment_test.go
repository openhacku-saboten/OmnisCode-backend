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

func TestCommentRepository_FindByID(t *testing.T) {
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

	if err := dbMap.Insert(&PostDTO{
		ID:        1,
		UserID:    "user-id",
		Title:     "test title",
		Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
		Language:  "Go",
		Content:   "Test code",
		Source:    "github.com",
		CreatedAt: time.Unix(100, 0),
		UpdatedAt: time.Unix(100, 0),
	}); err != nil {
		t.Fatal(err)
	}

	dbMap.AddTableWithName(CommentDTO{}, "comments").SetKeys(true, "id")
	truncateTable(t, dbMap, "comments")

	if err := dbMap.Insert(&CommentDTO{
		ID:        1,
		UserID:    "user-id",
		PostID:    1,
		Type:      "none",
		Content:   "type none",
		CreatedAt: time.Unix(100, 0),
		UpdatedAt: time.Unix(100, 0),
	}); err != nil {
		t.Fatal(err)
	}

	commentRepo := NewCommentRepository(dbMap)

	tests := []struct {
		name        string
		postID      int
		commentID   int
		wantComment *entity.Comment
		wantErr     error
	}{
		{
			name:      "正しくコメントを取得できる",
			postID:    1,
			commentID: 1,
			wantComment: &entity.Comment{
				ID:        1,
				UserID:    "user-id",
				PostID:    1,
				Type:      "none",
				Content:   "type none",
				CreatedAt: "1970-01-01T00:01:40+09:00",
				UpdatedAt: "1970-01-01T00:01:40+09:00",
			},
			wantErr: nil,
		},
		{
			name:        "PostIDが存在しなければErrNotFound",
			postID:      100,
			commentID:   1,
			wantComment: nil,
			wantErr:     entity.NewErrorNotFound("comment"),
		},
		{
			name:        "CommentIDが存在しなければErrNotFound",
			postID:      1,
			commentID:   100,
			wantComment: nil,
			wantErr:     entity.NewErrorNotFound("comment"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotComment, err := commentRepo.FindByID(context.Background(), tt.postID, tt.commentID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if diff := cmp.Diff(tt.wantComment, gotComment); diff != "" {
					t.Errorf("Data (-want +got) =\n%s\n", diff)
				}
			}
		})
	}
}

func TestCommentRepository_FindByPostID(t *testing.T) {
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
			CreatedAt: time.Unix(100, 0),
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
			CreatedAt: time.Unix(100, 0),
			UpdatedAt: time.Unix(100, 0),
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
			UserID:    "user-id",
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
			ctx := context.Background()
			gotComments, err := commentRepo.FindByPostID(ctx, tt.postID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				for idx := range gotComments {
					diff := cmp.Diff(tt.wantComments[idx], gotComments[idx], cmpopts.IgnoreFields(entity.Comment{}, "CreatedAt", "UpdatedAt"))
					if diff != "" {
						t.Errorf("Data (-want +got) =\n%s\n", diff)
					}
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

func TestCommentRepository_Insert(t *testing.T) {
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

	if err := dbMap.Insert(&PostDTO{
		ID:        1,
		UserID:    "user-id",
		Title:     "test title",
		Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
		Language:  "Go",
		Content:   "Test code",
		Source:    "github.com",
		CreatedAt: time.Unix(100, 0),
		UpdatedAt: time.Unix(100, 0),
	}); err != nil {
		t.Fatal(err)
	}

	dbMap.AddTableWithName(CommentDTO{}, "comments").SetKeys(true, "id")
	truncateTable(t, dbMap, "comments")

	if err := dbMap.Insert(&CommentDTO{
		ID:        1,
		UserID:    "user-id",
		PostID:    1,
		Type:      "none",
		Content:   "type none",
		CreatedAt: time.Unix(100, 0),
		UpdatedAt: time.Unix(100, 0),
	}); err != nil {
		t.Fatal(err)
	}

	commentRepo := NewCommentRepository(dbMap)

	tests := []struct {
		name          string
		comment       *entity.Comment
		wantErr       error
		wantCommentID int
	}{
		{
			name: "正しくコメントを作成できる",
			comment: &entity.Comment{
				UserID:  "user-id",
				PostID:  1,
				Type:    "none",
				Content: "type none",
			},
			wantErr:       nil,
			wantCommentID: 2,
		},
		{
			name: "IDが重複していてもAUTO_INCREMENTしてくれる",
			comment: &entity.Comment{
				ID:      1,
				UserID:  "user-id",
				PostID:  1,
				Type:    "none",
				Content: "type none",
			},
			wantErr:       nil,
			wantCommentID: 3,
		},
		{
			name: "PostIDが存在しなければErrNotFound",
			comment: &entity.Comment{
				UserID:  "user-id",
				PostID:  100,
				Type:    "none",
				Content: "type none",
			},
			wantErr:       entity.NewErrorNotFound("post"),
			wantCommentID: 0,
		},
		{
			name: "UserIDが存在しなければErrNotFound",
			comment: &entity.Comment{
				UserID:  "non-existing-id",
				PostID:  1,
				Type:    "none",
				Content: "type none",
			},
			wantErr:       entity.NewErrorNotFound("user"),
			wantCommentID: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := commentRepo.Insert(ctx, tt.comment)

			if tt.comment.ID != tt.wantCommentID {
				t.Errorf("gotCommentID = %d, want = %d", tt.comment.ID, tt.wantCommentID)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCommentRepository_Update(t *testing.T) {
	dbMap, err := NewDB()
	if err != nil {
		t.Fatalf(err.Error())
	}

	dbMap.AddTableWithName(UserDTO{}, "users")
	truncateTable(t, dbMap, "users")

	preparedUsers := []*UserDTO{
		{
			ID:        "user-id",
			Name:      "test user",
			Profile:   "test profile",
			TwitterID: "twitter",
		},
		{
			ID:        "user-id2",
			Name:      "test user2",
			Profile:   "test profile2",
			TwitterID: "twitter2",
		},
	}

	for _, user := range preparedUsers {
		if err := dbMap.Insert(user); err != nil {
			t.Fatal(err)
		}
	}

	dbMap.AddTableWithName(PostDTO{}, "posts").SetKeys(true, "id")
	dbMap.AddTableWithName(PostInsertDTO{}, "posts").SetKeys(true, "id")
	truncateTable(t, dbMap, "posts")

	validPosts := []*PostInsertDTO{
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
			UserID:   "user-id2",
			Title:    "test title",
			Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
			Language: "Go",
			Content:  "Test code",
			Source:   "github.com",
		},
		{
			ID:       3,
			UserID:   "user-id",
			Title:    "test title",
			Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
			Language: "Go",
			Content:  "Test code",
			Source:   "github.com",
		},
	}
	// デフォルトの投稿追加
	for _, post := range validPosts {
		if err := dbMap.Insert(post); err != nil {
			t.Fatal(err)
		}
	}

	dbMap.AddTableWithName(CommentDTO{}, "comments").SetKeys(true, "id")
	dbMap.AddTableWithName(CommentInsertDTO{}, "comments").SetKeys(true, "id")
	truncateTable(t, dbMap, "comments")

	if err := dbMap.Insert(&CommentInsertDTO{
		ID:      1,
		UserID:  "user-id",
		PostID:  1,
		Type:    "none",
		Content: "type none",
	}); err != nil {
		t.Fatal(err)
	}

	commentRepo := NewCommentRepository(dbMap)

	tests := []struct {
		name    string
		comment *entity.Comment
		wantErr error
	}{
		{
			name: "正しくコメントを更新できる",
			comment: &entity.Comment{
				ID:      1,
				UserID:  "user-id",
				PostID:  1,
				Type:    "none",
				Content: "type none",
			},
			wantErr: nil,
		},
		{
			name: `存在しないユーザが更新するとErrNotFoundにしたいが、
			これはUseCaseで実現するのでひとまずErrIsNotAuthorにしておく`,
			comment: &entity.Comment{
				ID:      1,
				UserID:  "user-id100",
				PostID:  1,
				Type:    "none",
				Content: "type none",
			},
			wantErr: entity.ErrIsNotAuthor,
		},
		{
			name: "投稿元のユーザ以外が更新するとErrIsNotAuthor",
			comment: &entity.Comment{
				ID:      1,
				UserID:  "user-id2",
				PostID:  1,
				Type:    "none",
				Content: "type none",
			},
			wantErr: entity.ErrIsNotAuthor,
		},
		{
			name: `PostIDが存在しなければErrNotFound`,
			comment: &entity.Comment{
				UserID:  "user-id",
				PostID:  100,
				Type:    "none",
				Content: "type none",
			},
			wantErr: entity.NewErrorNotFound("commnent"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := commentRepo.Update(context.Background(), tt.comment)

			errNF := &entity.ErrNotFound{}
			if errors.As(err, errNF) {
				if errors.As(tt.wantErr, errNF) {
					return
				}
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCommentRepository_Delete(t *testing.T) {
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

	if err := dbMap.Insert(&PostDTO{
		ID:        1,
		UserID:    "user-id",
		Title:     "test title",
		Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
		Language:  "Go",
		Content:   "Test code",
		Source:    "github.com",
		CreatedAt: time.Unix(100, 0),
		UpdatedAt: time.Unix(100, 0),
	}); err != nil {
		t.Fatal(err)
	}

	dbMap.AddTableWithName(CommentDTO{}, "comments").SetKeys(true, "id")
	truncateTable(t, dbMap, "comments")

	commentDTOs := []*CommentDTO{
		{
			ID:        1,
			UserID:    "user-id",
			PostID:    1,
			Type:      "none",
			Content:   "type none",
			CreatedAt: time.Unix(100, 0),
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
			CreatedAt: time.Unix(100, 0),
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
		name    string
		comment *entity.Comment
		wantErr error
	}{
		{
			name: "正しくコメントを削除できる",
			comment: &entity.Comment{
				ID:     1,
				PostID: 1,
				UserID: "user-id",
			},
			wantErr: nil,
		},
		{
			name: "コメントが存在しないならErrNotFound",
			comment: &entity.Comment{
				ID:     100,
				PostID: 1,
				UserID: "user-id",
			},
			wantErr: entity.NewErrorNotFound("comment"),
		},
		{
			name: "ユーザーが違うならErrIsNotAuthor",
			comment: &entity.Comment{
				ID:     2,
				PostID: 1,
				UserID: "other-user-id",
			},
			wantErr: entity.ErrIsNotAuthor,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := commentRepo.Delete(context.Background(), tt.comment)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
		})
	}
}
