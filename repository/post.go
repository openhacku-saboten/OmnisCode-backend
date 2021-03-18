//go:generate mockgen -source=$GOFILE -package=mock -destination=../usecase/mock/mock_$GOFILE

package repository

import (
	"context"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

// Post は永続化と再構築のためのインタフェースです
type Post interface {
	GetAll(ctx context.Context) ([]*entity.Post, error)
}
