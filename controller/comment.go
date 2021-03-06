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

// CommentController は コメントに関するハンドラに対してHTTPリクエストとして
// 送られたデータを入力として、ユースケースに伝えるまでを責務とするコントローラです
type CommentController struct {
	uc *usecase.CommentUseCase
}

// NewCommentController はCommentControllerのポインタを生成する関数です
func NewCommentController(uc *usecase.CommentUseCase) *CommentController {
	return &CommentController{uc: uc}
}

// Get は GET /post/{postID}/comment/{commentID} のHandler
func (ctrl *CommentController) Get(c echo.Context) error {
	logger := log.New()
	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	commentID, err := strconv.Atoi(c.Param("commentID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	comment, err := ctrl.uc.Get(c.Request().Context(), postID, commentID)

	if err != nil {
		errNF := &entity.ErrNotFound{}
		if errors.As(err, errNF) {
			return echo.NewHTTPError(http.StatusNotFound, errNF.Error())
		}

		logger.Errorf("Unexpected error GET /post/{postID}/comment/{commentID}: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, comment)
}

// GetByPostID は GET /post/{postID}/comment のHandler
func (ctrl *CommentController) GetByPostID(c echo.Context) error {
	logger := log.New()
	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	comments, err := ctrl.uc.GetByPostID(c.Request().Context(), postID)

	if err != nil {
		errNF := &entity.ErrNotFound{}
		if errors.As(err, errNF) {
			return echo.NewHTTPError(http.StatusNotFound, errNF.Error())
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
		logger.Errorf("Failed type assertion of userID: %#v", c.Get("userID"))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	comment.UserID = userID

	if err := ctrl.uc.Create(c.Request().Context(), comment); err != nil {
		if errors.Is(err, entity.ErrCannotCommit) {
			return echo.NewHTTPError(http.StatusForbidden, entity.ErrCannotCommit.Error())
		}
		errNF := &entity.ErrNotFound{}
		if errors.As(err, errNF) {
			return echo.NewHTTPError(http.StatusNotFound, errNF.Error())
		}
		logger.Errorf("error POST /post/{postID}/comment: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, comment)
}

// Update は PUT /post/{postID}/comment/{commentID} のHandler
func (ctrl *CommentController) Update(c echo.Context) error {
	logger := log.New()

	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	commentID, err := strconv.Atoi(c.Param("commentID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	comment := &entity.Comment{}
	if err := c.Bind(comment); err != nil {
		logger.Info(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	comment.ID = commentID
	comment.PostID = postID

	var ok bool
	comment.UserID, ok = c.Get("userID").(string)
	if !ok {
		logger.Errorf("Failed type assertion of userID: %#v", c.Get("userID"))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if err := ctrl.uc.Update(c.Request().Context(), comment); err != nil {
		if errors.Is(err, entity.ErrCannotCommit) {
			// コミットできない場合は、StatusForbidden
			return echo.NewHTTPError(http.StatusForbidden, entity.ErrCannotCommit.Error())
		}
		errNF := &entity.ErrNotFound{}
		if errors.As(err, errNF) {
			return echo.NewHTTPError(http.StatusNotFound, errNF.Error())
		}
		logger.Errorf("error POST /post/{postID}/comment: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

// Delete は DELETE /post/{postID}/comment/{commentID} のHandler
func (ctrl *CommentController) Delete(c echo.Context) error {
	logger := log.New()

	var err error
	comment := &entity.Comment{}
	comment.ID, err = strconv.Atoi(c.Param("commentID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	comment.PostID, err = strconv.Atoi(c.Param("postID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	var ok bool
	comment.UserID, ok = c.Get("userID").(string)
	if !ok {
		logger.Errorf("Failed type assertion of userID: %#v", c.Get("userID"))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if err := ctrl.uc.Delete(c.Request().Context(), comment); err != nil {
		if errors.Is(err, entity.ErrIsNotAuthor) {
			return echo.NewHTTPError(http.StatusForbidden, entity.ErrIsNotAuthor.Error())
		}
		errNF := &entity.ErrNotFound{}
		if errors.As(err, errNF) {
			return echo.NewHTTPError(http.StatusNotFound, errNF.Error())
		}
		logger.Errorf("error POST /post/{postID}/comment: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
