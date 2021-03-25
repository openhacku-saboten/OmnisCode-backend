package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase/mock"
)

func TestPost_Get_All_With_Mock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validPosts := []*entity.Post{
		{
			ID:        1,
			UserID:    "testID",
			Title:     "test title",
			Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
			Language:  "Go",
			Content:   "Test code",
			Source:    "github.com",
			CreatedAt: "2021-03-23T11:42:56+09:00",
			UpdatedAt: "2021-03-23T11:42:56+09:00",
		},
		{
			ID:        2,
			UserID:    "testID",
			Title:     "test title",
			Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
			Language:  "Go",
			Content:   "Test code",
			Source:    "github.com",
			CreatedAt: "2021-03-23T11:42:56+09:00",
			UpdatedAt: "2021-03-23T11:42:56+09:00",
		},
	}

	ctx := context.Background()
	postMock := mock.NewMockPost(ctrl)
	postMock.EXPECT().GetAll(ctx).Return(validPosts, nil)
	userMock := mock.NewMockPost(ctrl)

	sut := NewPostUsecase(postMock, userMock)
	posts, err := sut.postRepo.GetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for idx := range posts {
		if diff := cmp.Diff(posts[idx], validPosts[idx]); diff != "" {
			t.Fatalf("GetAll: %s", diff)
		}
	}
}

func TestPost_Get_With_Mock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validPost := &entity.Post{
		ID:        1,
		UserID:    "testID",
		Title:     "test title",
		Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
		Language:  "Go",
		Content:   "Test code",
		Source:    "github.com",
		CreatedAt: "2021-03-23T11:42:56+09:00",
		UpdatedAt: "2021-03-23T11:42:56+09:00",
	}

	ctx := context.Background()
	postMock := mock.NewMockPost(ctrl)
	postMock.EXPECT().FindByID(ctx, 1).Return(validPost, nil)
	userMock := mock.NewMockUser(ctrl)

	sut := NewPostUsecase(postMock, userMock)
	post, err := sut.postRepo.FindByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(post, validPost); diff != "" {
		t.Fatalf("FindByID: %s", diff)
	}
}

func TestPost_Create_With_Mock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validPost := &entity.Post{
		ID:        0,
		UserID:    "testID",
		Title:     "test title",
		Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
		Language:  "Go",
		Content:   "Test code",
		Source:    "github.com",
		CreatedAt: "2021-03-23T11:42:56+09:00",
		UpdatedAt: "2021-03-23T11:42:56+09:00",
	}

	ctx := context.Background()
	postMock := mock.NewMockPost(ctrl)
	postMock.EXPECT().Insert(ctx, validPost).Return(nil)
	userMock := mock.NewMockUser(ctrl)

	sut := NewPostUsecase(postMock, userMock)
	if err := sut.Create(ctx, validPost); err != nil {
		t.Fatal(err)
	}
}

func TestPost_Update_With_Mock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validPost := &entity.Post{
		ID:        0,
		UserID:    "testID",
		Title:     "test title",
		Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
		Language:  "Go",
		Content:   "Test code",
		Source:    "github.com",
		CreatedAt: "2021-03-23T11:42:56+09:00",
		UpdatedAt: "2021-03-23T11:42:56+09:00",
	}

	ctx := context.Background()
	postMock := mock.NewMockPost(ctrl)
	postMock.EXPECT().Update(ctx, validPost).Return(nil)
	userMock := mock.NewMockUser(ctrl)
	userMock.EXPECT().FindByID(validPost.UserID).Return(nil, nil)

	sut := NewPostUsecase(postMock, userMock)
	if err := sut.Update(ctx, validPost); err != nil {
		t.Fatal(err)
	}
}
