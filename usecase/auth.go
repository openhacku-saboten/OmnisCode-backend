package usecase

import (
	"context"

	"github.com/openhacku-saboten/OmnisCode-backend/repository"
)

// AuthUseCase は認証に関するユースケースです
type AuthUseCase struct {
	authRepo repository.Auth
}

// NewAuthUseCase はAuthUseCaseのポインタを生成する関数です
func NewAuthUseCase(authRepo repository.Auth) *AuthUseCase {
	return &AuthUseCase{authRepo: authRepo}
}

// Authenticate は認証を行い、userIDを取得します
func (a *AuthUseCase) Authenticate(ctx context.Context, token string) (uid string, err error) {
	uid, err = a.authRepo.Authenticate(ctx, token)
	return
}
