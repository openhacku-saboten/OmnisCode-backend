package repository

import "github.com/openhacku-saboten/OmnisCode-backend/domain/entity"

type User interface {
	FindByID(uid string) (user *entity.User, err error)
}
