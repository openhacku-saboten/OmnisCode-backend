package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/log"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase"
)

// PostController は 投稿に関するハンドラに対してHTTPリクエストとして
// 送られたデータを入力として、ユースケースに伝えるまでを責務とするコントローラです
type PostController struct {
	uc *usecase.PostUsecase
}

// NewPostController はPostControllerのポインタを生成する関数です
func NewPostController(uc *usecase.PostUsecase) *PostController {
	return &PostController{uc: uc}
}

// GetAll は GET /postのためのハンドラです
func (ctrl *PostController) GetAll(c echo.Context) error {
	logger := log.New()

	posts, err := ctrl.uc.GetAll(c.Request().Context())
	if err != nil {
		errNF := &entity.ErrNotFound{}
		if errors.As(err, errNF) {
			return echo.NewHTTPError(http.StatusNotFound, errNF.Error())
		}

		logger.Errorf("error GET /post: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, posts)
}

// Get は GET /post/{postID}のハンドラです
func (ctrl *PostController) Get(c echo.Context) error {
	logger := log.New()

	postID := c.Param("postID")
	if len(postID) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		// 数字ではない場合はエラー
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	ctx := c.Request().Context()
	post, err := ctrl.uc.Get(ctx, postIDInt)

	if err != nil {
		errNF := &entity.ErrNotFound{}
		if errors.As(err, errNF) {
			logger.Error(entity.NewErrorNotFound("post").Error())
			return echo.NewHTTPError(http.StatusNotFound, errNF.Error())
		}

		logger.Errorf("unexpected error GET /post/{postID}: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, post)
}

// Create は POST /postのハンドラです
func (ctrl *PostController) Create(c echo.Context) error {
	logger := log.New()

	post := &entity.Post{}
	if err := c.Bind(post); err != nil {
		logger.Errorf("failed c.Bind: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	var ok bool
	post.UserID, ok = c.Get("userID").(string)
	if !ok {
		logger.Errorf("Failed type assertion of userID: %#v", c.Get("userID"))
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	ctx := c.Request().Context()
	if err := ctrl.uc.Create(ctx, post); err != nil {
		logger.Errorf("error POST /post: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

// Update は Post /post/{postID}のハンドラです
func (ctrl *PostController) Update(c echo.Context) error {
	logger := log.New()

	post := &entity.Post{}
	if err := c.Bind(post); err != nil {
		logger.Errorf("failed c.Bind: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	var ok bool
	post.UserID, ok = c.Get("userID").(string)
	if !ok {
		logger.Errorf("Failed type assertion of userID: %#v", c.Get("userID"))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	ctx := c.Request().Context()
	if err := ctrl.uc.Update(ctx, post); err != nil {
		logger.Errorf("error POST /post: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
