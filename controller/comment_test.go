package controller

import (
	"encoding/json"
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

func TestCommentController_GetByPostID(t *testing.T) {
	tests := []struct {
		name               string
		postID             string
		prepareMockComment func(comment *mock.MockComment)
		wantErr            bool
		wantCode           int
		wantBody           string
	}{
		{
			name:   "正しくコメントを取得できる",
			postID: "1",
			prepareMockComment: func(comment *mock.MockComment) {
				comment.EXPECT().FindByPostID(1).Return(
					[]*entity.Comment{
						{
							ID:        1,
							UserID:    "userid1",
							PostID:    1,
							Type:      "highlight",
							Content:   "content1",
							FirstLine: 10,
							LastLine:  12,
							Code:      "",
							CreatedAt: "1970-01-01T09:01:40+09:00",
							UpdatedAt: "1970-01-01T09:01:40+09:00",
						},
						{
							ID:        2,
							UserID:    "userid2",
							PostID:    1,
							Type:      "commit",
							Content:   "content2",
							FirstLine: 0,
							LastLine:  0,
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
			name:   "postIDが空ならBadRequest",
			postID: "",
			prepareMockComment: func(comment *mock.MockComment) {
			},
			wantErr:  true,
			wantCode: 400,
			wantBody: "",
		},
		{
			name:   "postIDが数値でないならBadRequest",
			postID: "a",
			prepareMockComment: func(comment *mock.MockComment) {
			},
			wantErr:  true,
			wantCode: 400,
			wantBody: "",
		},
		{
			name:   "取得したコメント数が0ならErrNotFound",
			postID: "100",
			prepareMockComment: func(comment *mock.MockComment) {
				comment.EXPECT().FindByPostID(100).Return(
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
			c.SetParamNames("postID")
			c.SetParamValues(tt.postID)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commentRepo := mock.NewMockComment(ctrl)
			tt.prepareMockComment(commentRepo)
			postRepo := mock.NewMockPost(ctrl)
			userRepo := mock.NewMockUser(ctrl)

			con := NewCommentController(usecase.NewCommentUseCase(commentRepo, postRepo, userRepo))
			err := con.GetByPostID(c)

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

func TestCommentController_Create(t *testing.T) {
	tests := []struct {
		name               string
		postID             string
		userID             string
		body               string
		prepareMockComment func(comment *mock.MockComment)
		prepareMockPost    func(post *mock.MockPost)
		wantErr            bool
		wantCode           int
	}{
		{
			name:   "正しくコメントを作成できる",
			postID: "1",
			userID: "user-id",
			body: `{
				"type": "highlight",
				"content": "content1",
				"first_line": 10,
				"last_line": 12
			}`,
			prepareMockComment: func(comment *mock.MockComment) {
				comment.EXPECT().Insert(
					&entity.Comment{
						UserID:    "user-id",
						PostID:    1,
						Type:      "highlight",
						Content:   "content1",
						FirstLine: 10,
						LastLine:  12,
					}).Return(nil)
			},
			prepareMockPost: func(post *mock.MockPost) {
				post.EXPECT().FindByID(gomock.Any(), 1).Return(
					&entity.Post{
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
			wantCode: 201,
		},
		{
			name:   "Postのオーナー以外によるcommitならErrCannotCommit",
			postID: "1",
			userID: "user-id",
			body: `{
				"type": "commit",
				"content": "content1",
				"code":"hello"
			}`,
			prepareMockComment: func(comment *mock.MockComment) {
			},
			prepareMockPost: func(post *mock.MockPost) {
				post.EXPECT().FindByID(gomock.Any(), 1).Return(
					&entity.Post{
						ID:        1,
						UserID:    "other-user-id",
						Title:     "test title",
						Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
						Language:  "Go",
						Content:   "Test code",
						Source:    "github.com",
						CreatedAt: "2021-03-23T11:42:56+09:00",
						UpdatedAt: "2021-03-23T11:42:56+09:00",
					}, nil)
			},
			wantErr:  true,
			wantCode: 400,
		},
		{
			name:   "存在しないPostIDならErrNotFound",
			postID: "100",
			userID: "user-id",
			body: `{
				"type": "highlight",
				"content": "content1",
				"first_line": 10,
				"last_line": 12
			}`,
			prepareMockComment: func(comment *mock.MockComment) {
			},
			prepareMockPost: func(post *mock.MockPost) {
				post.EXPECT().FindByID(gomock.Any(), 100).Return(nil, entity.NewErrorNotFound("post"))
			},
			wantErr:  true,
			wantCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest("POST", "/", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("postID")
			c.SetParamValues(tt.postID)
			c.Set("userID", tt.userID)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commentRepo := mock.NewMockComment(ctrl)
			tt.prepareMockComment(commentRepo)
			postRepo := mock.NewMockPost(ctrl)
			tt.prepareMockPost(postRepo)
			userRepo := mock.NewMockUser(ctrl)

			con := NewCommentController(usecase.NewCommentUseCase(commentRepo, postRepo, userRepo))
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

func TestCommentController_Get(t *testing.T) {
	tests := []struct {
		name               string
		postID             string
		commentID          string
		prepareMockComment func(comment *mock.MockComment)
		wantErr            bool
		wantCode           int
		wantBody           string
	}{
		{
			name:      "正しくコメントを取得できる",
			postID:    "1",
			commentID: "1",
			prepareMockComment: func(comment *mock.MockComment) {
				comment.EXPECT().FindByID(1, 1).Return(
					&entity.Comment{
						ID:        1,
						UserID:    "userid1",
						PostID:    1,
						Type:      "highlight",
						Content:   "content1",
						FirstLine: 10,
						LastLine:  12,
						CreatedAt: "1970-01-01T09:01:40+09:00",
						UpdatedAt: "1970-01-01T09:01:40+09:00",
					}, nil)
			},
			wantErr:  false,
			wantCode: 200,
			wantBody: `{
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
			}`,
		},
		{
			name:      "postIDが空ならBadRequest",
			postID:    "",
			commentID: "1",
			prepareMockComment: func(comment *mock.MockComment) {
			},
			wantErr:  true,
			wantCode: 400,
			wantBody: "",
		},
		{
			name:      "postIDが数値でないならBadRequest",
			postID:    "a",
			commentID: "1",
			prepareMockComment: func(comment *mock.MockComment) {
			},
			wantErr:  true,
			wantCode: 400,
			wantBody: "",
		},
		{
			name:      "commentIDが空ならBadRequest",
			postID:    "1",
			commentID: "",
			prepareMockComment: func(comment *mock.MockComment) {
			},
			wantErr:  true,
			wantCode: 400,
			wantBody: "",
		},
		{
			name:      "commentIDが数値でないならBadRequest",
			postID:    "1",
			commentID: "a",
			prepareMockComment: func(comment *mock.MockComment) {
			},
			wantErr:  true,
			wantCode: 400,
			wantBody: "",
		},
		{
			name:      "commentが存在しないならNotFound",
			postID:    "1",
			commentID: "1",
			prepareMockComment: func(comment *mock.MockComment) {
				comment.EXPECT().FindByID(1, 1).Return(
					nil, entity.NewErrorNotFound("comment"))
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
			c.SetParamNames("postID", "commentID")
			c.SetParamValues(tt.postID, tt.commentID)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commentRepo := mock.NewMockComment(ctrl)
			tt.prepareMockComment(commentRepo)
			postRepo := mock.NewMockPost(ctrl)
			userRepo := mock.NewMockUser(ctrl)

			con := NewCommentController(usecase.NewCommentUseCase(commentRepo, postRepo, userRepo))
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
