package infra

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/service"
)

func TestPostRepository_GetAll(t *testing.T) {
	dbMap, err := NewDB()
	if err != nil {
		t.Fatalf(err.Error())
	}
	dbMap.AddTableWithName(UserDTO{}, "users")
	truncateTable(t, dbMap, "users")

	dbMap.AddTableWithName(PostDTO{}, "posts").SetKeys(true, "id")
	dbMap.AddTableWithName(PostInsertDTO{}, "posts").SetKeys(true, "id")

	if err := dbMap.Insert(&UserDTO{
		ID:        "user-id",
		Name:      "test user",
		Profile:   "test profile",
		TwitterID: "twitter",
	}); err != nil {
		t.Fatal(err)
	}

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
			UserID:   "user-id",
			Title:    "test title",
			Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
			Language: "Go",
			Content:  "Test code",
			Source:   "github.com",
		},
	}

	postRepo := NewPostRepository(dbMap)
	wantPosts := []*entity.Post{}

	for _, validPost := range validPosts {
		wantPosts = append(wantPosts, &entity.Post{
			ID:       validPost.ID,
			UserID:   validPost.UserID,
			Title:    validPost.Title,
			Code:     validPost.Code,
			Language: validPost.Language,
			Content:  validPost.Content,
			Source:   validPost.Source,
		})
	}

	tests := []struct {
		name      string
		posts     []*entity.Post
		wantPosts []*entity.Post
		wantErr   error
	}{
		{
			name:      "正しく全ての投稿を取得できる",
			posts:     wantPosts,
			wantPosts: wantPosts,
			wantErr:   nil,
		},
		{
			name:      "投稿が存在しなければNotFound",
			posts:     nil,
			wantPosts: nil,
			wantErr:   entity.NewErrorNotFound("post"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			// 初期化
			truncateTable(t, dbMap, "posts")

			for _, validPost := range tt.posts {
				// デフォルトの投稿追加
				if err := postRepo.Insert(ctx, validPost); err != nil {
					t.Fatal(err)
				}
			}
			posts, err := postRepo.GetAll(ctx)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				for idx := range posts {
					diff := cmp.Diff(tt.wantPosts[idx], posts[idx], cmpopts.IgnoreFields(entity.Post{}, "CreatedAt", "UpdatedAt"))
					if diff != "" {
						t.Errorf("Data (-want +got) =\n%s\n", diff)
					}
				}
			}
		})
	}

}

func TestPostRepository_FindByID(t *testing.T) {
	dbMap, err := NewDB()
	if err != nil {
		t.Fatalf(err.Error())
	}
	dbMap.AddTableWithName(PostDTO{}, "posts").SetKeys(true, "id")
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

	dbMap.AddTableWithName(PostDTO{}, "posts").SetKeys(true, "id")
	dbMap.AddTableWithName(PostInsertDTO{}, "posts").SetKeys(true, "id")
	truncateTable(t, dbMap, "posts")

	validPost := &PostInsertDTO{
		ID:       1,
		UserID:   "user-id",
		Title:    "test title",
		Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
		Language: "Go",
		Content:  "Test code",
		Source:   "github.com",
	}
	// デフォルトの投稿追加
	if err := dbMap.Insert(validPost); err != nil {
		t.Fatal(err)
	}

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
			postRepo := NewPostRepository(dbMap)
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

func TestUserRepository_FindByUserID(t *testing.T) {
	dbMap, err := NewDB()
	if err != nil {
		t.Fatalf(err.Error())
	}
	dbMap.AddTableWithName(UserDTO{}, "users")
	truncateTable(t, dbMap, "users")

	dbMap.AddTableWithName(PostDTO{}, "posts")
	dbMap.AddTableWithName(PostInsertDTO{}, "posts")

	validUsers := []*UserDTO{
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

	for _, user := range validUsers {
		if err := dbMap.Insert(user); err != nil {
			t.Fatal(err)
		}
	}

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

	var wantPosts []*entity.Post
	for _, validPost := range validPosts {
		wantPosts = append(wantPosts, &entity.Post{
			ID:       validPost.ID,
			UserID:   validPost.UserID,
			Title:    validPost.Title,
			Code:     validPost.Code,
			Language: validPost.Language,
			Content:  validPost.Content,
			Source:   validPost.Source,
		})
	}

	tests := []struct {
		name      string
		userID    string
		posts     []*entity.Post
		wantPosts []*entity.Post
		wantErr   error
	}{
		{
			name:   "正しく全ての投稿を取得できる",
			userID: "user-id",
			posts:  wantPosts,
			wantPosts: []*entity.Post{
				wantPosts[0],
				wantPosts[2],
			},
			wantErr: nil,
		},
		{
			name:      "投稿が存在しなければNotFound",
			userID:    "user-id3",
			posts:     wantPosts,
			wantPosts: nil,
			wantErr:   entity.NewErrorNotFound("post"),
		},
	}

	postRepo := NewPostRepository(dbMap)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			// 初期化
			truncateTable(t, dbMap, "posts")

			for _, validPost := range tt.posts {
				// デフォルトの投稿追加
				if err := postRepo.Insert(ctx, validPost); err != nil {
					t.Fatal(err)
				}
			}
			posts, err := postRepo.FindByUserID(ctx, tt.userID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				for idx := range posts {
					diff := cmp.Diff(tt.wantPosts[idx], posts[idx], cmpopts.IgnoreFields(entity.Post{}, "CreatedAt", "UpdatedAt"))
					if diff != "" {
						t.Errorf("Data (-want +got) =\n%s\n", diff)
					}
				}
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

	dbMap.AddTableWithName(PostDTO{}, "posts").SetKeys(true, "id")
	dbMap.AddTableWithName(PostInsertDTO{}, "posts").SetKeys(true, "id")
	truncateTable(t, dbMap, "posts")

	// デフォルトの投稿追加
	if err := dbMap.Insert(&PostInsertDTO{
		ID:       1,
		UserID:   "user-id",
		Title:    "test title",
		Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
		Language: "Go",
		Content:  "Test code",
		Source:   "github.com",
	}); err != nil {
		t.Fatal(err)
	}

	postRepo := NewPostRepository(dbMap)

	tests := []struct {
		name       string
		post       *entity.Post
		wantErr    error
		wantPostID int
	}{
		{
			name: "正常に追加できる",
			post: &entity.Post{
				UserID:   "user-id",
				Title:    "test title",
				Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language: "Go",
				Content:  "Test code",
				Source:   "github.com",
			},
			wantErr:    nil,
			wantPostID: 2,
		},
		{
			name: "存在しないユーザで登録するとエラー",
			post: &entity.Post{
				UserID:   "user-id2",
				Title:    "test title",
				Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language: "Go",
				Content:  "Test code",
				Source:   "github.com",
			},
			wantErr:    errors.New("unexisted user"),
			wantPostID: 0,
		},
		{
			name: "重複したpostIDで登録しても、auto incrementが働いてエラーは発生しない",
			post: &entity.Post{
				ID:       1,
				UserID:   "user-id",
				Title:    "test title",
				Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language: "Go",
				Content:  "Test code",
				Source:   "github.com",
			},
			wantErr:    nil,
			wantPostID: 4,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := postRepo.Insert(ctx, tt.post)

			if tt.post.ID != tt.wantPostID {
				t.Errorf("gotPostID = %d, want = %d", tt.post.ID, tt.wantPostID)
			}

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

func TestPostRepository_Update(t *testing.T) {
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

	postRepo := NewPostRepository(dbMap)

	tests := []struct {
		name    string
		post    *entity.Post
		wantErr error
	}{
		{
			name: "正常に更新できる",
			post: &entity.Post{
				ID:       1,
				UserID:   "user-id",
				Title:    "test title2",
				Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language: "Go",
				Content:  "Test code",
				Source:   "github.com",
			},
			wantErr: nil,
		},
		{
			name: `存在しないユーザで登録するとErrNotFoundにしたいが、
			Gorpは検知してくれないのでUsecaseでUserRepositoryを用いて
			存在証明をするので、エラーは投稿元のユーザ以外が更新すると
			エラーの時と同じerrIsNotAuthor`,
			post: &entity.Post{
				ID:       3,
				UserID:   "user-id3",
				Title:    "test title",
				Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language: "Go",
				Content:  "Test code",
				Source:   "github.com",
			},
			wantErr: entity.ErrIsNotAuthor,
		},
		{
			name: "投稿元のユーザ以外が更新するとエラー",
			post: &entity.Post{
				ID:       1,
				UserID:   "user-id2",
				Title:    "test title2",
				Code:     "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				Language: "Go",
				Content:  "Test code",
				Source:   "github.com",
			},
			wantErr: entity.ErrIsNotAuthor,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := postRepo.Update(ctx, tt.post)

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

func TestPostRepository_Delete(t *testing.T) {
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

	postRepo := NewPostRepository(dbMap)

	tests := []struct {
		name    string
		post    *entity.Post
		wantErr error
	}{
		{
			name: "正常に削除できる",
			post: &entity.Post{
				ID:     1,
				UserID: "user-id",
			},
			wantErr: nil,
		},
		{
			name: `存在しないPostならErrNotFound`,
			post: &entity.Post{
				ID:     100,
				UserID: "user-id",
			},
			wantErr: entity.NewErrorNotFound("post"),
		},
		{
			name: "投稿元のユーザ以外が削除するとエラー",
			post: &entity.Post{
				ID:     2,
				UserID: "user-id100",
			},
			wantErr: entity.ErrIsNotAuthor,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := postRepo.Delete(ctx, tt.post)

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
