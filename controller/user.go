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
