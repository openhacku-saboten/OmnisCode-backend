package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/log"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase"
)

type UserController struct {
	uc *usecase.UserUseCase
}

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
		if errors.Is(err, entity.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, entity.ErrUserNotFound.Error())
		}

		logger.Errorf("Unexpected error GET/user/{userID}: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, user)
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
		errNF := &ErrorNotFound{}
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

	if err := ctrl.uc.Create(user); err != nil {
		if errors.Is(err, entity.ErrDuplicatedUser) {
			return echo.NewHTTPError(http.StatusBadRequest, entity.ErrDuplicatedUser.Error())
		}
		if errors.Is(err, entity.ErrDuplicatedTwitterID) {
			return echo.NewHTTPError(http.StatusBadRequest, entity.ErrDuplicatedTwitterID.Error())
		}
		if errors.Is(err, entity.ErrEmptyUserName) {
			return echo.NewHTTPError(http.StatusBadRequest, entity.ErrEmptyUserName.Error())
		}
		logger.Errorf("Unexpected error POST/user: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
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

	if err := ctrl.uc.Update(user); err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, entity.ErrUserNotFound.Error())
		}
		if errors.Is(err, entity.ErrEmptyUserName) {
			return echo.NewHTTPError(http.StatusBadRequest, entity.ErrEmptyUserName.Error())
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

	return nil
}
