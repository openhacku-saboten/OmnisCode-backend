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

type CommentController struct {
	uc *usecase.CommentUseCase
}

func NewCommentController(uc *usecase.CommentUseCase) *CommentController {
	return &CommentController{uc: uc}
}

// GetByPostID は GET /post/{postID}/comment のHandler
func (ctrl *CommentController) GetByPostID(c echo.Context) error {
	logger := log.New()
	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	comments, err := ctrl.uc.GetByPostID(postID)

	if err != nil {
		if errors.As(err, &entity.ErrNotFound{}) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		logger.Errorf("Unexpected error GET /post/{postID}/comment: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, comments)
}

// Create は POST /post/{postID}/comment のHandler
func (ctrl *CommentController) Create(c echo.Context) error {
	logger := log.New()

	comment := &entity.Comment{}
	if err := c.Bind(comment); err != nil {
		logger.Info(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil {
		logger.Info(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	comment.PostID = postID

	userID, ok := c.Get("userID").(string)
	if !ok {
		if errors.Is(err, entity.ErrCannotCommit) {
			return echo.NewHTTPError(http.StatusBadRequest, entity.ErrCannotCommit.Error())
		}
		logger.Errorf("Failed type assertion of userID: %#v", c.Get("userID"))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	comment.UserID = userID

	if err := ctrl.uc.Create(c.Request().Context(), comment); err != nil {
		logger.Errorf("error POST /post/{postID}/comment: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}
