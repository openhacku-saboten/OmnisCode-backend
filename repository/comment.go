//go:generate mockgen -source=$GOFILE -destination=../usecase/mock/mock_$GOFILE -package=mock

package repository

import (
	"context"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

type Comment interface {
	GetByPostID(postid int) (comments []*entity.Comment, err error)
	FindByUserID(ctx context.Context, uid string) ([]*entity.Comment, error)
}
