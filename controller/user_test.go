package controller

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/labstack/echo/v4"
// 	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
// 	"github.com/openhacku-saboten/OmnisCode-backend/usecase"
// 	"github.com/openhacku-saboten/OmnisCode-backend/usecase/mock"
// )

// func TestUserController_Get(t *testing.T) {
// 	type fields struct {
// 		uc *usecase.UserUseCase
// 	}
// 	type args struct {
// 		c echo.Context
// 	}
// 	tests := []struct {
// 		name string
// 		userID string
// 		prepareMockUser func(u *mock.MockUser)
// 		wantErr  bool
// 		wantCode int
// 		wantData *entity.User
// 	}{
// 		{
// 			name:"正しくユーザーを取得できる",
// 			userID: "user-id",
// 			prepareMockUser: func(u *mock.MockUser) {
// 				u.EXPECT().FindByID("user-id").Return(
// 					entity.NewUser("user-id", "username", "profile", "@twitter"),
// 					nil,
// 				)
// 			},
// 			wantErr: false,
// 			wantCode: http.StatusOK,
// 			wantData: entity.NewUser("user-id", "username", "profile", "@twitter"),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			e := echo.New()
// 			req := httptest.NewRequest("GET", "/", nil)
// 			rec := httptest.NewRecorder()
// 			c := e.NewContext(req, rec)
// 			c.SetParamNames("userID")
// 			c.SetParamValues(tt.userID)

// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()
// 			userRepo := mock.NewMockUser(ctrl)
// 			tt.prepareMockUser(userRepo)

// 			m := NewUserController(usecase.NewUserUseCase(userRepo))
// 			err := m.Get(c)

// 			// ここから
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
// 			}

// 			if er, ok := err.(*echo.HTTPError); (ok && er.Code != tt.wantCode) || (!ok && rec.Code != tt.wantCode) {
// 				t.Errorf("code = %d, want = %d", rec.Code, tt.wantCode)
// 			}
// 		})
// 	}
// }
