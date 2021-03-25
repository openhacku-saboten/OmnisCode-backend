//go:generate mockgen -source=$GOFILE -destination=../infra/mock/mock_$GOFILE -package=mock

package repository

import (
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

// User はユーザに関する永続化と再構築のためのリポジトリです
type User interface {
	FindByID(uid string) (user *entity.User, err error)
	Insert(user *entity.User) error
	Update(user *entity.User) error
}
