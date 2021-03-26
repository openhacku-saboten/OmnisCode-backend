package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/echo/v4"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase/mock"
)

func TestPostController_GetAll(t *testing.T) {
	tests := []struct {
		name            string
		prepareMockPost func(ctx context.Context, post *mock.MockPost)
		wantErr         bool
		wantCode        int
		wantBody        string
	}{
		{
			name: "正しく投稿を取得できる",
			prepareMockPost: func(ctx context.Context, post *mock.MockPost) {
				post.EXPECT().GetAll(ctx).Return([]*entity.Post{
					{
						ID:        1,
						UserID:    "user-id",
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
						UserID:    "user-id",
						Title:     "test title",
						Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
						Language:  "Go",
						Content:   "Test code",
						Source:    "github.com",
						CreatedAt: "2021-03-23T11:42:56+09:00",
						UpdatedAt: "2021-03-23T11:42:56+09:00",
					},
				}, nil)
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			wantBody: `[{"id":1,"user_id":"user-id","title":"test title","code":"package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}","language":"Go","content":"Test code","source":"github.com","created_at":"2021-03-23T11:42:56+09:00","updated_at":"2021-03-23T11:42:56+09:00"},{"id":2,"user_id":"user-id","title":"test title","code":"package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}","language":"Go","content":"Test code","source":"github.com","created_at":"2021-03-23T11:42:56+09:00","updated_at":"2021-03-23T11:42:56+09:00"}]
`,
		},
		{
			name: "1つも投稿が存在しないならErrUserNotFound",
			prepareMockPost: func(ctx context.Context, post *mock.MockPost) {
				post.EXPECT().GetAll(ctx).Return(nil, entity.NewErrorNotFound("post"))
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
			wantBody: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ctx := c.Request().Context()
			postRepo := mock.NewMockPost(ctrl)
			tt.prepareMockPost(ctx, postRepo)
			userRepo := mock.NewMockUser(ctrl)

			con := NewPostController(usecase.NewPostUsecase(postRepo, userRepo))
			err := con.GetAll(c)

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

			if got := rec.Body.String(); got != tt.wantBody {
				t.Errorf("\nwant: %s, \nbut: %s", tt.wantBody, got)
			}
		})
	}
}

func TestPostController_Get(t *testing.T) {
	tests := []struct {
		name            string
		postID          string
		prepareMockPost func(ctx context.Context, post *mock.MockPost)
		wantErr         bool
		wantCode        int
	}{
		{
			name:   "正しく投稿を取得できる",
			postID: "1",
			prepareMockPost: func(ctx context.Context, post *mock.MockPost) {
				post.EXPECT().FindByID(ctx, 1).Return(&entity.Post{
					ID:        1,
					UserID:    "user-id",
					Title:     "test title",
					Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
					Language:  "Go",
					Content:   "Test code",
					Source:    "github.com",
					CreatedAt: "2021-03-23T11:42:56+09:00",
					UpdatedAt: "2021-03-23T11:42:56+09:00",
				}, nil)
			},
			wantErr:  false,
			wantCode: http.StatusOK,
		},
		{
			name:   "存在しない投稿IDならErrUserNotFound",
			postID: "0",
			prepareMockPost: func(ctx context.Context, post *mock.MockPost) {
				post.EXPECT().FindByID(ctx, 0).Return(&entity.Post{
					ID:        1,
					UserID:    "user-id",
					Title:     "test title",
					Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
					Language:  "Go",
					Content:   "Test code",
					Source:    "github.com",
					CreatedAt: "2021-03-23T11:42:56+09:00",
					UpdatedAt: "2021-03-23T11:42:56+09:00",
				}, entity.NewErrorNotFound("post"))
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("postID")
			c.SetParamValues(tt.postID)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ctx := c.Request().Context()
			postRepo := mock.NewMockPost(ctrl)
			tt.prepareMockPost(ctx, postRepo)
			userRepo := mock.NewMockUser(ctrl)

			con := NewPostController(usecase.NewPostUsecase(postRepo, userRepo))
			err := con.Get(c)

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

func TestPostController_Create(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		body            string
		prepareMockPost func(ctx context.Context, post *mock.MockPost)
		wantErr         bool
		wantCode        int
		wantBody        string
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
					UserID:    "user-id",
					Title:     "test title",
					Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
					Language:  "Go",
					Content:   "Test code",
					Source:    "github.com",
					CreatedAt: "2021-03-23T11:42:56+09:00",
					UpdatedAt: "2021-03-23T11:42:56+09:00",
				}).DoAndReturn(func(ctx context.Context, user *entity.Post) error {
					user.ID = 1
					return nil
				})
			},
			wantErr:  false,
			wantCode: 201,
			wantBody: `{
				"id": 1,
				"user_id":"user-id",
				"title":"test title",
				"code":"package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				"language":"Go",
				"content":"Test code",
				"source":"github.com",
				"created_at":"2021-03-23T11:42:56+09:00",
				"updated_at":"2021-03-23T11:42:56+09:00"
				}`,
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
			ctx := context.Background()
			tt.prepareMockPost(ctx, postRepo)
			userRepo := mock.NewMockUser(ctrl)

			con := NewPostController(usecase.NewPostUsecase(postRepo, userRepo))
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

			if !tt.wantErr {
				var gotBody, wantBody map[string]interface{}
				if err = json.Unmarshal(rec.Body.Bytes(), &gotBody); err != nil {
					t.Fatal(err)
				}
				if err = json.Unmarshal([]byte(tt.wantBody), &wantBody); err != nil {
					t.Fatal(err)
				}

				if diff := cmp.Diff(wantBody, gotBody); diff != "" {
					t.Errorf("body (-want +got) =\n%s\n", diff)
				}
			}
		})
	}
}

func TestPostController_Update(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		body            string
		prepareMockPost func(ctx context.Context, post *mock.MockPost)
		wantErr         bool
		wantCode        int
	}{
		{
			name:   "正しく投稿を更新できる",
			userID: "user-id",
			body: `{
				"id": 1,
				"title":"test title",
				"code":"package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				"language":"Go",
				"content":"Test code",
				"source":"github.com",
				"created_at":"2021-03-23T11:42:56+09:00",
				"updated_at":"2021-03-23T11:42:56+09:00"
				}`,
			prepareMockPost: func(ctx context.Context, post *mock.MockPost) {
				post.EXPECT().Update(ctx, &entity.Post{
					ID:        1,
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
			wantCode: 200,
		},
		{
			name:   "不正なBodyならBadRequest",
			userID: "user-id",
			body: `{
				"aaaa":"test title",
				}`,
			prepareMockPost: func(ctx context.Context, post *mock.MockPost) {},
			wantErr:         true,
			wantCode:        http.StatusBadRequest,
		},
		{
			name:            "bodyがJSON形式でないならBadRequest",
			userID:          "user-id",
			body:            `aaaaa`,
			prepareMockPost: func(ctx context.Context, post *mock.MockPost) {},
			wantErr:         true,
			wantCode:        http.StatusBadRequest,
		},
		{
			name:   "存在しないポストならばErrIsNotAuthorでForbidden",
			userID: "user-id",
			body: `{
				"id": 100,
				"title":"test title",
				"code":"package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				"language":"Go",
				"content":"Test code",
				"source":"github.com",
				"created_at":"2021-03-23T11:42:56+09:00",
				"updated_at":"2021-03-23T11:42:56+09:00"
				}`,
			prepareMockPost: func(ctx context.Context, post *mock.MockPost) {
				post.EXPECT().Update(ctx, &entity.Post{
					ID:        100,
					UserID:    "user-id",
					Title:     "test title",
					Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
					Language:  "Go",
					Content:   "Test code",
					Source:    "github.com",
					CreatedAt: "2021-03-23T11:42:56+09:00",
					UpdatedAt: "2021-03-23T11:42:56+09:00",
				}).Return(entity.ErrIsNotAuthor)
			},
			wantErr:  true,
			wantCode: http.StatusForbidden,
		},
		{
			name:   "存在しないユーザならばErrNotFoundでForbidden",
			userID: "user-id2002",
			body: `{
				"id": 1,
				"title":"test title",
				"code":"package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
				"language":"Go",
				"content":"Test code",
				"source":"github.com",
				"created_at":"2021-03-23T11:42:56+09:00",
				"updated_at":"2021-03-23T11:42:56+09:00"
				}`,
			prepareMockPost: func(ctx context.Context, post *mock.MockPost) {
				post.EXPECT().Update(ctx, &entity.Post{
					ID:        1,
					UserID:    "user-id2002",
					Title:     "test title",
					Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
					Language:  "Go",
					Content:   "Test code",
					Source:    "github.com",
					CreatedAt: "2021-03-23T11:42:56+09:00",
					UpdatedAt: "2021-03-23T11:42:56+09:00",
				}).Return(entity.NewErrorNotFound("user"))
			},
			wantErr:  true,
			wantCode: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest("PUT", "/", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("userID", tt.userID)

			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			postRepo := mock.NewMockPost(ctrl)
			tt.prepareMockPost(ctx, postRepo)
			userRepo := mock.NewMockUser(ctrl)

			con := NewPostController(usecase.NewPostUsecase(postRepo, userRepo))
			err := con.Update(c)

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
