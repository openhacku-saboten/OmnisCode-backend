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

	authUseCase := usecase.NewAuthUseCase(authRepo)
	authMiddleware := controller.NewAuthMiddleware(authUseCase)

	userUseCase := usecase.NewUserUseCase(userRepo, authRepo)
	userController := controller.NewUserController(userUseCase)
	_ = userController

	postUsecase := usecase.NewPostUsecase(postRepo)
	postController := controller.NewPostController(postUsecase)

	e := echo.New()
	v1 := e.Group("/api/v1")

	user := v1.Group("/user")
	user.GET("/:userID", userController.Get)
	user.POST("", userController.Create, authMiddleware.Authenticate)

	post := v1.Group("/post")
	post.POST("", postController.Create, authMiddleware.Authenticate)

	if err := e.Start(fmt.Sprintf(":%s", config.Port())); err != nil {
		logger.Infof("shutting down the server with error' %s", err.Error())
		os.Exit(1)
	}
}
