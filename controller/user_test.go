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
				user.EXPECT().FindByID("user-id").Return(
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
				user.EXPECT().FindByID("invalid-user-id").Return(
					nil,
					entity.ErrUserNotFound,
				)
			},
			prepareMockAuth: func(auth *mock.MockAuth) {},
			wantErr:         true,
			wantCode:        404,
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

			con := NewUserController(usecase.NewUserUseCase(userRepo, authRepo))
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
					entity.NewUser("user-id", "username", "profile", "twitter", ""),
				).Return(nil)
			},
			wantErr:  false,
			wantCode: 200,
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
					entity.NewUser("user-id", "username", "profile", "twitter", ""),
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
			name:   "userIDが重複しているならBadRequest",
			userID: "user-id",
			body: `{
				"name":"username",
				"profile":"profile",
				"twitter_id":"twitter"
			}`,
			prepareMockUser: func(user *mock.MockUser) {
				user.EXPECT().Insert(
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

			con := NewUserController(usecase.NewUserUseCase(userRepo, authRepo))
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
