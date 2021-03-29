//go:generate mockgen -source=$GOFILE -destination=../infra/mock/mock_$GOFILE -package=mock

package repository

import (
	"context"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

// Auth はFirebase Authentication関連の操作を表すインターフェース
type Auth interface {
	Authenticate(ctx context.Context, token string) (uid string, err error)
	GetIconURL(ctx context.Context, uid string) (iconURL string, err error)
	Delete(ctx context.Context, user *entity.User) error
}
