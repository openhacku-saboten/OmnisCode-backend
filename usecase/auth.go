package usecase

import "github.com/openhacku-saboten/OmnisCode-backend/repository"

type AuthUseCase struct {
	authRepo repository.Auth
}

func NewAuthUseCase(authRepo repository.Auth) *AuthUseCase {
	return &AuthUseCase{authRepo: authRepo}
}

func (a *AuthUseCase) Authenticate(token string) (uid string, err error) {
	uid, err = a.authRepo.Authenticate(token)
	return
}
