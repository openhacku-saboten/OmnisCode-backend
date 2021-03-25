package main

import (
	"context"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/openhacku-saboten/OmnisCode-backend/config"
	"github.com/openhacku-saboten/OmnisCode-backend/controller"
	"github.com/openhacku-saboten/OmnisCode-backend/infra"
	"github.com/openhacku-saboten/OmnisCode-backend/log"
	"github.com/openhacku-saboten/OmnisCode-backend/usecase"
)

func main() {
	logger := log.New()

	dbMap, err := infra.NewDB()
	if err != nil {
		logger.Errorf("failed NewDB: %s", err.Error())
		os.Exit(1)
	}
	defer func() {
		err := dbMap.Db.Close()
		if err != nil {
			logger.Errorf("failed to close DB: %s", err.Error())
		}
	}()

	ctx := context.Background()
	firebase, err := infra.NewFirebase(ctx)
	if err != nil {
		logger.Errorf("failed NewFirebase: %s", err.Error())
		os.Exit(1)
	}

	authRepo := infra.NewAuthRepository(firebase)
	userRepo := infra.NewUserRepository(dbMap)
	postRepo := infra.NewPostRepository(dbMap)
	commentRepo := infra.NewCommentRepository(dbMap)

	authUseCase := usecase.NewAuthUseCase(authRepo)
	authMiddleware := controller.NewAuthMiddleware(authUseCase)

	userUseCase := usecase.NewUserUseCase(userRepo, authRepo, postRepo, commentRepo)
	userController := controller.NewUserController(userUseCase)

	postUsecase := usecase.NewPostUsecase(postRepo, userRepo)
	postController := controller.NewPostController(postUsecase)

	commentUseCase := usecase.NewCommentUseCase(commentRepo, postRepo)
	commentController := controller.NewCommentController(commentUseCase)

	e := echo.New()
	v1 := e.Group("/api/v1")

	user := v1.Group("/user")
	user.GET("/:userID", userController.Get)
	user.POST("", userController.Create, authMiddleware.Authenticate)
	user.PUT("", userController.Update, authMiddleware.Authenticate)
	user.GET("/:userID/post", userController.GetPosts)
	user.GET("/:userID/comment", userController.GetComments)

	post := v1.Group("/post")
	post.GET("", postController.GetAll) // 記事の閲覧はログインの必要なし
	post.POST("", postController.Create, authMiddleware.Authenticate)
	post.GET("/:postID", postController.Get)
	post.PUT("/:postID", postController.Update, authMiddleware.Authenticate)

	comment := v1.Group("/post/:postID/comment")
	comment.GET("", commentController.GetByPostID)
	comment.POST("", commentController.Create, authMiddleware.Authenticate)
	comment.GET("/:commentID", commentController.Get)

	if err := e.Start(fmt.Sprintf(":%s", config.Port())); err != nil {
		logger.Infof("shutting down the server with error' %s", err.Error())
		os.Exit(1)
	}
}
