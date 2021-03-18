package usecase

import "github.com/openhacku-saboten/OmnisCode-backend/repository"

type AuthUseCase struct {
	authRepo repository.Auth
}

func NewAuthUseCase(authRepo repository.Auth) *AuthUseCase {
	return &AuthUseCase{authRepo: authRepo}
}

func (u *AuthUseCase) Authenticate(token string) (uid string, err error) {
	uid, err = u.authRepo.Authenticate(token)
	return
}
