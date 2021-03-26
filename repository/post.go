//go:generate mockgen -source=$GOFILE -destination=../infra/mock/mock_$GOFILE -package=mock

package repository

import (
	"context"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

// Post は投稿に関する永続化と再構成のためのリポジトリです
type Post interface {
	GetAll(ctx context.Context) ([]*entity.Post, error)
	FindByID(ctx context.Context, postID int) (*entity.Post, error)
	FindByUserID(ctx context.Context, uid string) ([]*entity.Post, error)
	Insert(ctx context.Context, post *entity.Post) error
	Update(ctx context.Context, post *entity.Post) error
	Delete(ctx context.Context, userID string, postID int) error
}
