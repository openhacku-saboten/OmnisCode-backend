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

func TestUserController_Get(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		prepareMockUser func(user *mock.MockUser)
		prepareMockAuth func(auth *mock.MockAuth)
		wantErr         bool
		wantCode        int
		wantBody        map[string]interface{}
	}{
		{
			name:   "正しくユーザーを取得できる",
			userID: "user-id",
			prepareMockUser: func(user *mock.MockUser) {
				user.EXPECT().FindByID(gomock.Any(), "user-id").Return(
					entity.NewUser("user-id", "name", "profile", "twitter", ""),
					nil,
				)
			},
			prepareMockAuth: func(auth *mock.MockAuth) {
				auth.EXPECT().GetIconURL(gomock.Any(), "user-id").Return("icon-url", nil)
			},
			wantErr:  false,
			wantCode: 200,
			wantBody: map[string]interface{}{
				"id":         "user-id",
				"name":       "name",
				"profile":    "profile",
				"twitter_id": "twitter",
				"icon_url":   "icon-url",
			},
		},
		{
			name:   "存在しないユーザーIDならErrUserNotFound",
			userID: "invalid-user-id",
			prepareMockUser: func(user *mock.MockUser) {
				user.EXPECT().FindByID(gomock.Any(), "invalid-user-id").Return(
					nil,
					entity.ErrUserNotFound,
				)
			},
			prepareMockAuth: func(auth *mock.MockAuth) {},
			wantErr:         true,
			wantCode:        404,
			wantBody:        nil,
		},
		{
			name:   "ユーザーIDが空ならBadRequest",
			userID: "",
			prepareMockUser: func(user *mock.MockUser) {
			},
			prepareMockAuth: func(auth *mock.MockAuth) {},
			wantErr:         true,
			wantCode:        400,
			wantBody:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("userID")
			c.SetParamValues(tt.userID)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userRepo := mock.NewMockUser(ctrl)
			tt.prepareMockUser(userRepo)
			authRepo := mock.NewMockAuth(ctrl)
			tt.prepareMockAuth(authRepo)
			postRepo := mock.NewMockPost(ctrl)
			commentRepo := mock.NewMockComment(ctrl)

			con := NewUserController(usecase.NewUserUseCase(userRepo, authRepo, postRepo, commentRepo))
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

			if !tt.wantErr {
				var gotBody map[string]interface{}
				err = json.Unmarshal(rec.Body.Bytes(), &gotBody)
				if err != nil {
					t.Fatal(err)
				}

				if diff := cmp.Diff(tt.wantBody, gotBody); diff != "" {
					t.Errorf("body (-want +got) =\n%s\n", diff)
				}
			}
		})
	}
}

func TestCommentController_GetByUserID(t *testing.T) {
	tests := []struct {
		name               string
		userID             string
		prepareMockComment func(ctx context.Context, uid string, comment *mock.MockComment)
		wantErr            bool
		wantCode           int
		wantBody           string
	}{
		{
			name:   "正しくコメントを取得できる",
			userID: "userid1",
			prepareMockComment: func(ctx context.Context, uid string, comment *mock.MockComment) {
				comment.EXPECT().FindByUserID(ctx, uid).Return(
					[]*entity.Comment{
						{
							ID:        1,
							UserID:    "userid1",
							PostID:    1,
							Type:      "highlight",
							Content:   "content1",
							FirstLine: 10,
							LastLine:  12,
							CreatedAt: "1970-01-01T09:01:40+09:00",
							UpdatedAt: "1970-01-01T09:01:40+09:00",
						},
						{
							ID:        2,
							UserID:    "userid2",
							PostID:    1,
							Type:      "commit",
							Content:   "content2",
							Code:      "code2",
							CreatedAt: "1970-01-01T09:01:40+09:00",
							UpdatedAt: "1970-01-01T09:01:40+09:00",
						},
					},
					nil,
				)
			},
			wantErr:  false,
			wantCode: 200,
			wantBody: `[
				{
					"id": 1,
					"user_id": "userid1",
					"post_id": 1,
					"type": "highlight",
					"content": "content1",
					"first_line": 10,
					"last_line": 12,
					"code": "",
					"created_at": "1970-01-01T09:01:40+09:00",
					"updated_at": "1970-01-01T09:01:40+09:00"
				},
				{
					"id": 2,
					"user_id": "userid2",
					"post_id": 1,
					"type": "commit",
					"content": "content2",
					"first_line": 0,
					"last_line": 0,
					"code": "code2",
					"created_at": "1970-01-01T09:01:40+09:00",
					"updated_at": "1970-01-01T09:01:40+09:00"
				}
			]`,
		},
		{
			name:   "userIDが空ならBadRequest",
			userID: "",
			prepareMockComment: func(ctx context.Context, uid string, comment *mock.MockComment) {
			},
			wantErr:  true,
			wantCode: 400,
			wantBody: "",
		},
		{
			name:   "取得したコメント数が0ならErrNotFound",
			userID: "100",
			prepareMockComment: func(ctx context.Context, uid string, comment *mock.MockComment) {
				comment.EXPECT().FindByUserID(ctx, uid).Return(
					nil, entity.NewErrorNotFound("comment"),
				)
			},
			wantErr:  true,
			wantCode: 404,
			wantBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetParamNames("userID")
			c.SetParamValues(tt.userID)

			ctx := c.Request().Context()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userRepo := mock.NewMockUser(ctrl)
			authRepo := mock.NewMockAuth(ctrl)
			postRepo := mock.NewMockPost(ctrl)
			commentRepo := mock.NewMockComment(ctrl)
			tt.prepareMockComment(ctx, tt.userID, commentRepo)

			userCon := NewUserController(usecase.NewUserUseCase(userRepo, authRepo, postRepo, commentRepo))
			err := userCon.GetComments(c)

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
				var gotBody, wantBody []map[string]interface{}
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

func TestUserController_GetPosts(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		prepareMockPost func(ctx context.Context, uid string, post *mock.MockPost)
		wantErr         bool
		wantCode        int
		wantBody        string
	}{
		{
			name:   "正しく投稿を取得できる",
			userID: "user-id",
			prepareMockPost: func(ctx context.Context, uid string, post *mock.MockPost) {
				post.EXPECT().FindByUserID(ctx, uid).Return([]*entity.Post{
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
			name:   "1つも投稿が存在しないならErrUserNotFound",
			userID: "user-id2",
			prepareMockPost: func(ctx context.Context, uid string, post *mock.MockPost) {
				post.EXPECT().FindByUserID(ctx, uid).Return(nil, entity.NewErrorNotFound("post"))
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
			authRepo := mock.NewMockAuth(ctrl)
			userRepo := mock.NewMockUser(ctrl)
			postRepo := mock.NewMockPost(ctrl)
			tt.prepareMockPost(ctx, tt.userID, postRepo)
			commentRepo := mock.NewMockComment(ctrl)

			con := NewUserController(usecase.NewUserUseCase(userRepo, authRepo, postRepo, commentRepo))
			c.SetParamNames("userID")
			c.SetParamValues(tt.userID)
			err := con.GetPosts(c)

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
				var gotBody, wantBody []map[string]interface{}
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

func TestUserController_Create(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		body            string
		prepareMockUser func(user *mock.MockUser)
		prepareMockAuth func(auth *mock.MockAuth)
		wantErr         bool
		wantCode        int
	}{
		{
			name:   "正しくユーザーを作成できる",
			userID: "user-id",
			body: `{
				"name":"username",
				"profile":"profile",
				"twitter_id":"twitter"
			}`,
			prepareMockUser: func(user *mock.MockUser) {
				user.EXPECT().Insert(
					gomock.Any(),
					entity.NewUser("user-id", "username", "profile", "twitter", ""),
				).Return(nil)
			},
			wantErr:  false,
			wantCode: 201,
		},
		{
			name:   "TwitterIDに@が含まれていれば取り除いてユーザーを作成できる",
			userID: "user-id",
			body: `{
				"name":"username",
				"profile":"profile",
				"twitter_id":"@twitter"
			}`,
			prepareMockUser: func(user *mock.MockUser) {
				user.EXPECT().Insert(
					gomock.Any(),
					entity.NewUser("user-id", "username", "profile", "twitter", ""),
				).Return(nil)
			},
			wantErr:  false,
			wantCode: 201,
		},
		{
			name:   "不正なbodyならBadRequest",
			userID: "user-id",
			body: `{
				"aaa":"test"
			}`,
			prepareMockUser: func(user *mock.MockUser) {},
			wantErr:         true,
			wantCode:        400,
		},
		{
			name:            "bodyがJSON形式でないならBadRequest",
			userID:          "user-id",
			body:            `aaaaa`,
			prepareMockUser: func(user *mock.MockUser) {},
			wantErr:         true,
			wantCode:        400,
		},
		{
			name:   "userIDが重複しているならBadRequest",
			userID: "user-id",
			body: `{
				"name":"username",
				"profile":"profile",
				"twitter_id":"twitter"
			}`,
			prepareMockUser: func(user *mock.MockUser) {
				user.EXPECT().Insert(
					gomock.Any(),
					entity.NewUser("user-id", "username", "profile", "twitter", ""),
				).Return(entity.ErrDuplicatedUser)
			},
			wantErr:  true,
			wantCode: 400,
		},
		{
			name:   "TwitterIDが重複しているならBadRequest",
			userID: "user-id",
			body: `{
				"name":"username",
				"profile":"profile",
				"twitter_id":"twitter"
			}`,
			prepareMockUser: func(user *mock.MockUser) {
				user.EXPECT().Insert(
					gomock.Any(),
					entity.NewUser("user-id", "username", "profile", "twitter", ""),
				).Return(entity.ErrDuplicatedTwitterID)
			},
			wantErr:  true,
			wantCode: 400,
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
			userRepo := mock.NewMockUser(ctrl)
			tt.prepareMockUser(userRepo)
			authRepo := mock.NewMockAuth(ctrl)
			postRepo := mock.NewMockPost(ctrl)
			commentRepo := mock.NewMockComment(ctrl)

			con := NewUserController(usecase.NewUserUseCase(userRepo, authRepo, postRepo, commentRepo))
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

func TestUserController_Update(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		body            string
		prepareMockUser func(user *mock.MockUser)
		prepareMockAuth func(auth *mock.MockAuth)
		wantErr         bool
		wantCode        int
	}{
		{
			name:   "正しくユーザーを更新できる",
			userID: "user-id",
			body: `{
				"name":"newname",
				"profile":"newprofile",
				"twitter_id":"newtwitter"
			}`,
			prepareMockUser: func(user *mock.MockUser) {
				user.EXPECT().FindByID(gomock.Any(), "user-id").Return(
					entity.NewUser("user-id", "name", "profile", "twitter", ""),
					nil,
				)
				user.EXPECT().Update(
					gomock.Any(),
					entity.NewUser("user-id", "newname", "newprofile", "newtwitter", ""),
				).Return(nil)
			},
			wantErr:  false,
			wantCode: 200,
		},
		{
			name:   "TwitterIDに@が含まれていれば取り除いてユーザーを更新できる",
			userID: "user-id",
			body: `{
				"name":"newname",
				"profile":"newprofile",
				"twitter_id":"@newtwitter"
			}`,
			prepareMockUser: func(user *mock.MockUser) {
				user.EXPECT().FindByID(gomock.Any(), "user-id").Return(
					entity.NewUser("user-id", "name", "profile", "twitter", ""),
					nil,
				)
				user.EXPECT().Update(
					gomock.Any(),
					entity.NewUser("user-id", "newname", "newprofile", "newtwitter", ""),
				).Return(nil)
			},
			wantErr:  false,
			wantCode: 200,
		},
		{
			name:   "不正なbodyならBadRequest",
			userID: "user-id",
			body: `{
				"aaa":"test"
			}`,
			prepareMockUser: func(user *mock.MockUser) {},
			wantErr:         true,
			wantCode:        400,
		},
		{
			name:            "bodyがJSON形式でないならBadRequest",
			userID:          "user-id",
			body:            `aaaaa`,
			prepareMockUser: func(user *mock.MockUser) {},
			wantErr:         true,
			wantCode:        400,
		},
		{
			name:   "存在しないユーザーIDならErrUserNotFound",
			userID: "invalid-user-id",
			body: `{
				"name":"username",
				"profile":"profile",
				"twitter_id":"twitter"
			}`,
			prepareMockUser: func(user *mock.MockUser) {
				user.EXPECT().FindByID(gomock.Any(), "invalid-user-id").Return(
					nil,
					entity.ErrUserNotFound,
				)
			},
			prepareMockAuth: func(auth *mock.MockAuth) {},
			wantErr:         true,
			wantCode:        404,
		},
		{
			name:   "TwitterIDが重複しているならBadRequest",
			userID: "user-id",
			body: `{
				"name":"newname",
				"profile":"newprofile",
				"twitter_id":"newtwitter"
			}`,
			prepareMockUser: func(user *mock.MockUser) {
				user.EXPECT().FindByID(gomock.Any(), "user-id").Return(
					entity.NewUser("user-id", "name", "profile", "twitter", ""),
					nil,
				)
				user.EXPECT().Update(
					gomock.Any(),
					entity.NewUser("user-id", "newname", "newprofile", "newtwitter", ""),
				).Return(entity.NewErrorDuplicated("user TwitterID"))
			},
			wantErr:  true,
			wantCode: 400,
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

			userRepo := mock.NewMockUser(ctrl)
			tt.prepareMockUser(userRepo)
			authRepo := mock.NewMockAuth(ctrl)
			postRepo := mock.NewMockPost(ctrl)
			commentRepo := mock.NewMockComment(ctrl)

			con := NewUserController(usecase.NewUserUseCase(userRepo, authRepo, postRepo, commentRepo))
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
