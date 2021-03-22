package usecase

import (
	"context"

	"github.com/openhacku-saboten/OmnisCode-backend/repository"
)

type AuthUseCase struct {
	authRepo repository.Auth
}

func NewAuthUseCase(authRepo repository.Auth) *AuthUseCase {
	return &AuthUseCase{authRepo: authRepo}
}

func (a *AuthUseCase) Authenticate(ctx context.Context, token string) (uid string, err error) {
	uid, err = a.authRepo.Authenticate(ctx, token)
	return
}
