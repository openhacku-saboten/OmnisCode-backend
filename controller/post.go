package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/log"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase"
)

// PostController は投稿情報のHTTPリクエストをコントロールする構造体です
type PostController struct {
	uc *usecase.PostUsecase
}

// NewPostController はPostControllerのポインタを生成する関数です
func NewPostController(uc *usecase.PostUsecase) *PostController {
	return &PostController{uc: uc}
}

// Create は POST /postのハンドラです
func (ctrl *PostController) Create(c echo.Context) error {
	logger := log.New()

	post := &entity.Post{}

	logger.Info(c)
	if err := c.Bind(post); err != nil {
		logger.Infof("failed c.Bind: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	var ok bool
	post.UserID, ok = c.Get("userID").(string)
	if !ok {
		logger.Errorf("Failed type assertion of userID: %#v", c.Get("userID"))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if err := ctrl.uc.Create(c.Request().Context(), post); err != nil {
		logger.Errorf("error POST /post: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}
