package infra

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/google/go-cmp/cmp"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

func TestUserRepository_FindByID(t *testing.T) {
	dbMap, err := NewDB()
	if err != nil {
		t.Fatalf(err.Error())
	}
	dbMap.AddTableWithName(UserDTO{}, "users")
	truncateTable(t, dbMap, "users")
	if err := dbMap.Insert(&UserDTO{
		ID:        "existing-id",
		Name:      "existingUser",
		Profile:   "existing",
		TwitterID: "existing",
	}); err != nil {
		t.Fatal(err)
	}

	userRepo := NewUserRepository(dbMap)

	tests := []struct {
		name     string
		userID   string
		wantUser *entity.User
		wantErr  error
	}{
		{
			name:     "正しくユーザーを取得できる",
			userID:   "existing-id",
			wantUser: entity.NewUser("existing-id", "existingUser", "existing", "existing", ""),
			wantErr:  nil,
		},
		{
			name:     "存在しないユーザーの場合はErrNoRows",
			userID:   "not-existing-id",
			wantUser: nil,
			wantErr:  entity.ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := userRepo.FindByID(tt.userID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				diff := cmp.Diff(tt.wantUser, gotUser)
				if diff != "" {
					t.Errorf("Data (-want +got) =\n%s\n", diff)
				}
			}
		})
	}
}

func TestUserRepository_GetCommentsByID(t *testing.T) {
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
	dbMap.AddTableWithName(PostInsertDTO{}, "posts")
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

	dbMap.AddTableWithName(CommentDTO{}, "comments")
	dbMap.AddTableWithName(CommentInsertDTO{}, "comments")
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

	userRepo := NewUserRepository(dbMap)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			comments, err := userRepo.FindCommentsByID(ctx, tt.userID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				diff := cmp.Diff(tt.wantComments, comments)
				if diff != "" {
					t.Errorf("Data (-want +got) =\n%s\n", diff)
				}
			}
		})
	}
}

func TestUserRepository_Insert(t *testing.T) {
	dbMap, err := NewDB()
	if err != nil {
		t.Fatalf(err.Error())
	}
	dbMap.AddTableWithName(UserDTO{}, "users")
	truncateTable(t, dbMap, "users")
	if err := dbMap.Insert(&UserDTO{
		ID:        "existing-id",
		Name:      "existingUser",
		Profile:   "existing",
		TwitterID: "existing",
	}); err != nil {
		t.Fatal(err)
	}

	userRepo := NewUserRepository(dbMap)

	tests := []struct {
		name    string
		user    *entity.User
		wantErr error
	}{
		{
			name:    "すでに存在するユーザーIDならErrDuplicatedUser",
			user:    entity.NewUser("existing-id", "newUser", "new", "new", ""),
			wantErr: entity.ErrDuplicatedUser,
		},
		{
			name:    "すでに存在するTwitterIDならErrDuplicatedTwitterID",
			user:    entity.NewUser("new-id", "newUser", "new", "existing", ""),
			wantErr: entity.ErrDuplicatedTwitterID,
		},
		{
			name:    "正しくユーザーを作成できる",
			user:    entity.NewUser("new-id", "newUser", "new", "new", ""),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := userRepo.Insert(tt.user)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	dbMap, err := NewDB()
	if err != nil {
		t.Fatalf(err.Error())
	}
	dbMap.AddTableWithName(UserDTO{}, "users")
	truncateUser(t, dbMap)
	userDTOs := []*UserDTO{
		&UserDTO{
			ID:        "existing-id",
			Name:      "existingUser",
			Profile:   "existing",
			TwitterID: "existing",
		},
		&UserDTO{
			ID:        "existing-id2",
			Name:      "existingUser2",
			Profile:   "existing2",
			TwitterID: "existing2",
		},
	}
	for _, userDTO := range userDTOs {
		if err := dbMap.Insert(userDTO); err != nil {
			t.Fatal(err)
		}
	}

	userRepo := NewUserRepository(dbMap)

	tests := []struct {
		name    string
		user    *entity.User
		wantErr error
	}{
		{
			name:    "ユーザーIDが存在しないならDBを更新せずに終わる",
			user:    entity.NewUser("new-id", "newUser", "new", "new", ""),
			wantErr: nil,
		},
		{
			name:    "すでに存在するTwitterIDならErrDuplicatedTwitterID",
			user:    entity.NewUser("existing-id", "updateUser", "update", "existing2", ""),
			wantErr: entity.NewErrorDuplicated("user TwitterID"),
		},
		{
			name:    "フィールドに変更がなくても正しくユーザーを更新できる",
			user:    entity.NewUser("existing-id", "existingUser", "existing", "existing", ""),
			wantErr: nil,
		},
		{
			name:    "正しくユーザーを更新できる",
			user:    entity.NewUser("existing-id", "updateUser", "update", "update", ""),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := userRepo.Update(tt.user)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
		})
	}
}

func truncateUser(t *testing.T, dbMap *gorp.DbMap) {
	t.Helper()

	// databaseを初期化する
	if _, err := dbMap.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		t.Fatal(err)
	}
	// タイミングの問題でTruncateが失敗することがあるので成功するまで試みる
	for i := 0; i < 5; i++ {
		_, err := dbMap.Exec("TRUNCATE TABLE users")
		if err == nil {
			break
		}
		if i == 4 {
			t.Fatal(err)
		}
		time.Sleep(time.Second * 1)
	}
	if _, err := dbMap.Exec("SET FOREIGN_KEY_CHECKS = 1"); err != nil {
		t.Fatal(err)
	}
}
