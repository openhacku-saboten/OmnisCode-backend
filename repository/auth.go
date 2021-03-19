//go:generate mockgen -source=$GOFILE -destination=../usecase/mock/mock_$GOFILE -package=mock

package repository

import "context"

type Auth interface {
	Authenticate(ctx context.Context, token string) (uid string, err error)
}
