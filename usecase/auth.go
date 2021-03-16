package usecase

import "github.com/openhacku-saboten/OmnisCode-backend/repository"

type AuthUseCase struct {
	firebase repository.Firebase
}

func NewAuthUseCase(firebase repository.Firebase) *AuthUseCase {
	return &AuthUseCase{firebase: firebase}
}

func (u *AuthUseCase) Authenticate(token string) (uid string, err error) {
	uid, err = u.firebase.Authenticate(token)
	return
}
