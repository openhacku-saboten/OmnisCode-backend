package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/log"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase"
)

// UserController は ユーザに関するハンドラに対してHTTPリクエストとして
// 送られたデータを入力として、ユースケースに伝えるまでを責務とするコントローラです
type UserController struct {
	uc *usecase.UserUseCase
}

// NewUserController はUserControllerのポインタを生成する関数です
func NewUserController(uc *usecase.UserUseCase) *UserController {
	return &UserController{uc: uc}
}

// Get は GET /user/{userID} のHandler
func (ctrl *UserController) Get(c echo.Context) error {
	logger := log.New()
	userID := c.Param("userID")
	if len(userID) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	user, err := ctrl.uc.Get(c.Request().Context(), userID)

	if err != nil {
		errNF := &entity.ErrNotFound{}
		if errors.As(err, errNF) {
			return echo.NewHTTPError(http.StatusNotFound, errNF.Error())
		}

		logger.Errorf("Unexpected error GET/user/{userID}: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, user)
}

// GetPosts は  GET /user/{userID}/post のHandler
func (ctrl *UserController) GetPosts(c echo.Context) error {
	logger := log.New()

	userID := c.Param("userID")
	if len(userID) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	ctx := c.Request().Context()

	posts, err := ctrl.uc.GetPosts(ctx, userID)
	if err != nil {
		errNF := &entity.ErrNotFound{}
		if errors.As(err, errNF) {
			return echo.NewHTTPError(http.StatusNotFound, errNF.Error())
		}
		logger.Errorf("error GET /user/{userID}/post: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, posts)
}

// GetComments は GET /user/{userID}/comment
func (ctrl *UserController) GetComments(c echo.Context) error {
	logger := log.New()

	userID := c.Param("userID")
	if len(userID) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	ctx := c.Request().Context()
	comments, err := ctrl.uc.GetComments(ctx, userID)
	if err != nil {
		errNF := &entity.ErrNotFound{}
		if errors.As(err, errNF) {
			return echo.NewHTTPError(http.StatusNotFound, errNF.Error())
		}

		logger.Errorf("Unexpected error GET/user/{userID}/comment: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, comments)
}

// Create は POST /user のHandler
func (ctrl *UserController) Create(c echo.Context) error {
	logger := log.New()

	user := &entity.User{}
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	userID, ok := c.Get("userID").(string)
	if !ok {
		logger.Errorf("Failed type assertion of userID: %#v", c.Get("userID"))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	user.ID = userID

	if err := ctrl.uc.Create(c.Request().Context(), user); err != nil {
		errDup := &entity.ErrDuplicated{}
		if errors.As(err, errDup) {
			return echo.NewHTTPError(http.StatusBadRequest, errDup.Error())
		}
		errEmpty := &entity.ErrEmpty{}
		if errors.As(err, errEmpty) {
			return echo.NewHTTPError(http.StatusBadRequest, errEmpty.Error())
		}
		logger.Errorf("Unexpected error POST/user: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

// Update は PUT /user のHandler
func (ctrl *UserController) Update(c echo.Context) error {
	logger := log.New()

	user := &entity.User{}
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	userID, ok := c.Get("userID").(string)
	if !ok {
		logger.Errorf("Failed type assertion of userID: %#v", c.Get("userID"))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	user.ID = userID

	if err := ctrl.uc.Update(c.Request().Context(), user); err != nil {
		errNF := &entity.ErrNotFound{}
		if errors.As(err, errNF) {
			return echo.NewHTTPError(http.StatusNotFound, errNF.Error())
		}

		errEmpty := &entity.ErrEmpty{}
		if errors.As(err, errEmpty) {
			return echo.NewHTTPError(http.StatusBadRequest, errEmpty.Error())
		}
		errDup := &entity.ErrDuplicated{}
		if errors.As(err, errDup) {
			return echo.NewHTTPError(http.StatusBadRequest, errDup.Error())
		}
		errTL := &entity.ErrTooLong{}
		if errors.As(err, errTL) {
			return echo.NewHTTPError(http.StatusBadRequest, errTL.Error())
		}
		logger.Errorf("Unexpected error PUT/user: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

// Delete は DELETE /user/{userID} のHandler
func (ctrl *UserController) Delete(c echo.Context) error {
	logger := log.New()

	user := &entity.User{}
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	var ok bool
	user.ID, ok = c.Get("userID").(string)
	if !ok {
		logger.Errorf("Failed type assertion of userID: %#v", c.Get("userID"))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if err := ctrl.uc.Delete(c.Request().Context(), user); err != nil {
		errEmpty := &entity.ErrEmpty{}
		if errors.As(err, errEmpty) {
			return echo.NewHTTPError(http.StatusBadRequest, errEmpty.Error())
		}
		errNF := &entity.ErrNotFound{}
		if errors.As(err, errNF) {
			return echo.NewHTTPError(http.StatusNotFound, errNF.Error())
		}
		if errors.Is(err, entity.ErrIsNotAuthor) {
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		}
		logger.Errorf("Unexpected error PUT/user: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
