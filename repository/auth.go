//go:generate mockgen -source=$GOFILE -destination=../usecase/mock/mock_$GOFILE -package=mock

package repository

import "context"

// Auth はFirebase Authentication関連の操作を表すインターフェース
type Auth interface {
	Authenticate(ctx context.Context, token string) (uid string, err error)
}
