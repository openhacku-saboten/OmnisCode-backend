package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/openhacku-saboten/OmnisCode-backend/repository"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase"
)

type AuthMiddleware struct {
	uc *usecase.AuthUseCase
}

func NewAuthMiddleware(f repository.Firebase) *AuthMiddleware {
	return &AuthMiddleware{uc: usecase.NewAuthUseCase(f)}
}

func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		authScheme := "Bearer"

		// Tokenから"Bearer "を取り除く
		l := len(authScheme)
		if len(authHeader) <= l+1 || authHeader[:l] != authScheme {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid token")
		}
		token := authHeader[l+1:]

		userID, err := m.uc.Authenticate(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		c.Set("userID", userID)
		return next(c)
	}
}
