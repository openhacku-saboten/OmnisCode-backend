package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/openhacku-saboten/OmnisCode-backend/log"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase"
)

type Post struct {
	uc *usecase.Post
}

func NewPostController(uc *usecase.Post) *Post {
	return &Post{uc: uc}
}

func (p *Post) GetAll(c echo.Context) error {
	logger := log.New()
	ctx := c.Request().Context()

	posts, err := p.uc.GetAll(ctx)
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, &posts)
}
