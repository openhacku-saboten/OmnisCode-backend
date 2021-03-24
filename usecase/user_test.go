package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase/mock"
)

func TestUsecase_Get_Posts_By_ID(t *testing.T) {
	const (
		userID = "testID"
		token  = "test token"
	)
	validPosts := []*entity.Post{
		{
			ID:        1,
			UserID:    userID,
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
			UserID:    userID,
			Title:     "test title",
			Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
			Language:  "Go",
			Content:   "Test code",
			Source:    "github.com",
			CreatedAt: "2021-03-23T11:42:56+09:00",
			UpdatedAt: "2021-03-23T11:42:56+09:00",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	authMock := mock.NewMockAuth(ctrl)
	authMock.EXPECT().Authenticate(ctx, token).Return(userID, nil)
	userMock := mock.NewMockUser(ctrl)
	userMock.EXPECT().FindPostsByID(ctx, userID).Return(validPosts, nil)

	sut := NewUserUseCase(userMock, authMock)

	uid, err := sut.authRepo.Authenticate(ctx, token)
	if err != nil {
		t.Fatal(err)
	}
	posts, err := sut.userRepo.FindPostsByID(ctx, uid)
	if err != nil {
		t.Fatal(err)
	}

	for idx := range posts {
		if diff := cmp.Diff(posts[idx], validPosts[idx]); diff != "" {
			t.Fatalf("GetAll: %s", diff)
		}
	}
}
