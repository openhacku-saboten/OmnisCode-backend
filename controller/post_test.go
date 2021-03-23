package controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase/mock"
)

func TestPostController_Create(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		body            string
		prepareMockPost func(ctx context.Context, post *mock.MockPost)
		wantErr         bool
		wantCode        int
	}{
		{
			name:   "正しく投稿を作成できる",
			userID: "user-id",
			body: `{
				"title":"test title",
				"code":"package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				"language":"Go",
				"content":"Test code",
				"source":"github.com",
				"created_at":"2021-03-23T11:42:56+09:00",
				"updated_at":"2021-03-23T11:42:56+09:00"
				}`,
			prepareMockPost: func(ctx context.Context, post *mock.MockPost) {
				post.EXPECT().Insert(ctx, &entity.Post{
					ID:        0,
					UserID:    "user-id",
					Title:     "test title",
					Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
					Language:  "Go",
					Content:   "Test code",
					Source:    "github.com",
					CreatedAt: "2021-03-23T11:42:56+09:00",
					UpdatedAt: "2021-03-23T11:42:56+09:00",
				}).Return(nil)
			},
			wantErr:  false,
			wantCode: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest("POST", "/", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("userID", tt.userID)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			postRepo := mock.NewMockPost(ctrl)
			tt.prepareMockPost(context.Background(), postRepo)

			con := NewPostController(usecase.NewPostUsecase(postRepo))
			fmt.Println(c)
			err := con.Create(c)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
			}

			if he, ok := err.(*echo.HTTPError); ok {
				if he.Code != tt.wantCode {
					t.Errorf("code = %d, want = %d", he.Code, tt.wantCode)
				}
			} else {
				if rec.Code != tt.wantCode {
					t.Errorf("code = %d, want = %d", rec.Code, tt.wantCode)
				}
			}
		})
	}
}
