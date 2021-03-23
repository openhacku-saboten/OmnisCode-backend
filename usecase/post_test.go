package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase/mock"
)

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

	sut := NewPostUsecase(postMock)
	if err := sut.Create(ctx, validPost); err != nil {
		t.Fatal(err)
	}
}
