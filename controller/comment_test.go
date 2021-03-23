package controller

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

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
				comment.EXPECT().GetByPostID(1).Return(
					[]*entity.Comment{
						&entity.Comment{
							ID:        1,
							UserID:    "userid1",
							PostID:    1,
							Type:      "highlight",
							Content:   "content1",
							FirstLine: 10,
							LastLine:  12,
							CreatedAt: time.Unix(100, 0),
							UpdatedAt: time.Unix(100, 0),
						},
						&entity.Comment{
							ID:        2,
							UserID:    "userid2",
							PostID:    1,
							Type:      "commit",
							Content:   "content2",
							Code:      "code2",
							CreatedAt: time.Unix(100, 0),
							UpdatedAt: time.Unix(100, 0),
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

			con := NewCommentController(usecase.NewCommentUseCase(commentRepo))
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
