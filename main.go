package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/openhacku-saboten/OmnisCode-backend/config"
	"github.com/openhacku-saboten/OmnisCode-backend/controller"
	"github.com/openhacku-saboten/OmnisCode-backend/infra"
	"github.com/openhacku-saboten/OmnisCode-backend/log"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase"
)

func main() {
	logger := log.New()

	firebase, err := infra.NewFirebase()
	if err != nil {
		logger.Errorf("failed NewFirebase: %s", err.Error())
		os.Exit(1)
	}
	authRepo := infra.NewAuthRepository(firebase)
	authUseCase := usecase.NewAuthUseCase(authRepo)
	authMiddleware := controller.NewAuthMiddleware(authUseCase)

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		logger.Infof("Access from %s", c.Request().RemoteAddr)
		return c.String(http.StatusOK, "Hello, World!")
	})

	authed := e.Group("/secret", authMiddleware.Authenticate)

	authed.GET("", func(c echo.Context) error {
		logger.Infof("Authorized access from%s", c.Request().RemoteAddr)
		return c.String(http.StatusOK, c.Get("userID").(string))
	})

	if err := e.Start(fmt.Sprintf(":%s", config.Port())); err != nil {
		logger.Infof("shutting down the server with error' %s", err.Error())
		os.Exit(1)
	}
}
