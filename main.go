package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/openhacku-saboten/OmnisCode-backend/config"
	"github.com/openhacku-saboten/OmnisCode-backend/infra"
	"github.com/openhacku-saboten/OmnisCode-backend/log"
)

func main() {
	logger := log.New()

	dbMap, err := infra.NewDB()
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		err := dbMap.Db.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		logger.Infof("Access from %s", c.Request().RemoteAddr)
		return c.String(http.StatusOK, "Hello, World!")
	})

	if err := e.Start(fmt.Sprintf(":%s", config.Port())); err != nil {
		logger.Infof("shutting down the server with error' %s", err.Error())
		os.Exit(1)
	}
}
