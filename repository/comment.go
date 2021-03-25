//go:generate mockgen -source=$GOFILE -destination=../usecase/mock/mock_$GOFILE -package=mock

package repository

import (
	"context"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

// Comment はコメントに関する永続化と再構築のためのリポジトリです
type Comment interface {
	FindByUserID(ctx context.Context, uid string) ([]*entity.Comment, error)
	FindByPostID(postID int) (comments []*entity.Comment, err error)
	Insert(comment *entity.Comment) error
	FindByID(postID, commentID int) (comment *entity.Comment, err error)
	Delete(postID, commentID int) error
}
