package main

import (
	"fmt"
	"net/http"

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
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Port())))
}
