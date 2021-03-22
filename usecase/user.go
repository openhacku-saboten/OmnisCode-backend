package usecase

import (
	"context"
	"fmt"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/repository"
)

type UserUseCase struct {
	userRepo repository.User
	authRepo repository.Auth
}

func NewUserUseCase(user repository.User, auth repository.Auth) *UserUseCase {
	return &UserUseCase{userRepo: user, authRepo: auth}
}

func (u *UserUseCase) Get(ctx context.Context, uid string) (user *entity.User, err error) {
	user, err = u.userRepo.FindByID(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to Get User from DB: %w", err)
	}

	user.IconURL, err = u.authRepo.GetIconURL(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to Get User from Firebase: %w", err)
	}
	return
}
