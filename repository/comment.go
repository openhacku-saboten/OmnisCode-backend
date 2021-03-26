//go:generate mockgen -source=$GOFILE -destination=../infra/mock/mock_$GOFILE -package=mock

package repository

import (
	"context"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

// Comment はコメントに関する永続化と再構成のためのリポジトリです
type Comment interface {
	FindByID(ctx context.Context, postID, commentID int) (comment *entity.Comment, err error)
	FindByUserID(ctx context.Context, uid string) ([]*entity.Comment, error)
	FindByPostID(ctx context.Context, postID int) (comments []*entity.Comment, err error)
	Insert(ctx context.Context, comment *entity.Comment) error
	Update(ctx context.Context, comment *entity.Comment) error
}
