//go:generate mockgen -source=$GOFILE -destination=../usecase/mock/mock_$GOFILE -package=mock

package repository

import (
	"context"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

// Post は永続化と再構築のためのリポジトリです
type Post interface {
	GetAll(ctx context.Context) ([]*entity.Post, error)
	Insert(ctx context.Context, post *entity.Post) error
}
