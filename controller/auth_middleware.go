package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/openhacku-saboten/OmnisCode-backend/log"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase"
)

// AuthMiddleware はHTTPリクエストとして送られたデータを入力として
// ユースケースに伝えるまでを責務とするコントローラです
type AuthMiddleware struct {
	uc *usecase.AuthUseCase
}

// NewAuthMiddleware はAuthMiddlewareのポインタを生成する関数です
func NewAuthMiddleware(uc *usecase.AuthUseCase) *AuthMiddleware {
	return &AuthMiddleware{uc: uc}
}

// Authenticate はAuthorizationヘッダーを検証し，userIDをcontextにセットする(c.Get("userID")で取得可能)
func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	logger := log.New()
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		authScheme := "Bearer"

		// Tokenから"Bearer "を取り除く
		l := len(authScheme)
		if len(authHeader) <= l+1 || authHeader[:l] != authScheme {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid token")
		}
		token := authHeader[l+1:]

		userID, err := m.uc.Authenticate(c.Request().Context(), token)
		if err != nil {
			logger.Infof("Unauthorized: %v", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		c.Set("userID", userID)
		return next(c)
	}
}
