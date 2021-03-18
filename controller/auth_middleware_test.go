package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/echo/v4"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase/mock"
)

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name            string
		prepareRequest  func(req *http.Request)
		prepareFirebase func(f *mock.MockFirebase)
		next            echo.HandlerFunc
		wantErr         bool
		wantCode        int
	}{
		{
			name: "正しく認証できる",
			prepareRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer token")
			},
			prepareFirebase: func(f *mock.MockFirebase) {
				f.EXPECT().Authenticate("token").Return("currentUserID", nil)
			},
			next: func(c echo.Context) error {
				got, ok := c.Get("userID").(string)
				if !ok {
					t.Errorf("UserID not found in context")
				}
				want := "currentUserID"
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("Data (-want +got) =\n%s\n", diff)
				}
				return nil
			},
			wantErr:  false,
			wantCode: http.StatusOK,
		},
		{
			name: "HeaderがなければBadRequest",
			prepareRequest: func(req *http.Request) {
			},
			prepareFirebase: func(f *mock.MockFirebase) {
			},
			next:     nil,
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Headerの形式が不正ならBadRequest",
			prepareRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Token token")
			},
			prepareFirebase: func(f *mock.MockFirebase) {
			},
			next:     nil,
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "認証されていないTokenならUnauthorized",
			prepareRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer invalidToken")
			},
			prepareFirebase: func(f *mock.MockFirebase) {
				f.EXPECT().Authenticate("invalidToken").Return("error verifying ID token", errors.New("error verifying ID token"))
			},
			next:     nil,
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest("GET", "/", nil)
			tt.prepareRequest(req)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			firebase := mock.NewMockFirebase(ctrl)
			tt.prepareFirebase(firebase)

			m := NewAuthMiddleware(usecase.NewAuthUseCase(firebase))
			err := m.Authenticate(tt.next)(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if er, ok := err.(*echo.HTTPError); (ok && er.Code != tt.wantCode) || (!ok && rec.Code != tt.wantCode) {
				t.Errorf("code = %d, want = %d", rec.Code, tt.wantCode)
			}
		})
	}
}
