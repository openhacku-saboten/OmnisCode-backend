//go:generate mockgen -source=$GOFILE -destination=../usecase/mock/mock_$GOFILE -package=mock

package repository

import (
	"context"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

type User interface {
	FindByID(uid string) (user *entity.User, err error)
	FindCommentsByID(ctx context.Context, uid string) ([]*entity.Comment, error)
	Insert(user *entity.User) error
	Update(user *entity.User) error
}
