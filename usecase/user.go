package usecase

import (
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/repository"
)

type UserUseCase struct {
	User repository.User
}

func NewUserUseCase(user repository.User) *UserUseCase {
	return &UserUseCase{User: user}
}

func (u *UserUseCase) Get(uid string) (user *entity.User, err error) {
	user, err = u.User.FindByID(uid)
	return
}
