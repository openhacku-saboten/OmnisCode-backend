package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase/mock"
)

func TestPostController_Post(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	validPost := &entity.Post{
		ID:        0,
		UserID:    "testID",
		Title:     "test title",
		Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
		Language:  "Go",
		Content:   "Test code",
		Source:    "github.com",
		CreatedAt: time.Date(2021, time.March, 23, 11, 42, 56, 0, loc),
		UpdatedAt: time.Date(2021, time.March, 23, 11, 42, 56, 0, loc),
	}
	_ = validPost
	tests := []struct {
		name            string
		userID          string
		body            string
		prepareMockPost func(ctx context.Context, post *mock.MockPost)
		wantErr         bool
		wantCode        int
	}{
		{
			name:   "正しくユーザを作成できる",
			userID: "user-id",
			body: `{
				"Title":"test title",
				"Code":"package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				"Language":"Go",
				"Content":"Test code",
				"Source":"github.com",
				"CreatedAt":"2021-03-23T11:42:56 +09:00",
				"UpdatedAt":"2021-03-23T11:42:56+09:00"
				}`,
			prepareMockPost: func(ctx context.Context, post *mock.MockPost) {
				post.EXPECT().Store(ctx, &entity.Post{
					ID:        0,
					UserID:    "user-id",
					Title:     "test title",
					Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
					Language:  "Go",
					Content:   "Test code",
					Source:    "github.com",
					CreatedAt: time.Date(2021, time.March, 23, 11, 42, 56, 0, loc),
					UpdatedAt: time.Date(2021, time.March, 23, 11, 42, 56, 0, loc),
				}).Return(nil)
			},
			wantErr:  false,
			wantCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			c.SetParamNames("userID")
			c.SetParamValues(tt.userID)
			ctx := c.Request().Context()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			postRepo := mock.NewMockPost(ctrl)
			tt.prepareMockPost(ctx, postRepo)

			sut := NewPostController(usecase.NewPostUsecase(postRepo))
			err := sut.Create(c)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
			}

			if he, ok := err.(*echo.HTTPError); ok {
				if he.Code != tt.wantCode {
					t.Errorf("code = %d, want = %d", he.Code, tt.wantCode)
				}
			} else {
				if res.Code != tt.wantCode {
					t.Errorf("code = %d, want = %d", res.Code, tt.wantCode)
				}
			}
		})
	}
}
