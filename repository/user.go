//go:generate mockgen -source=$GOFILE -destination=../usecase/mock/mock_$GOFILE -package=mock

package repository

import (
	"context"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

// User はユーザに関する永続化と再構成のためのリポジトリです
type User interface {
	FindByID(ctx context.Context, uid string) (user *entity.User, err error)
	Insert(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
}
