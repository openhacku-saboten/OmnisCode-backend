//go:generate mockgen -source=$GOFILE -destination=../usecase/mock/mock_$GOFILE -package=mock

package repository

type Firebase interface {
	Authenticate(token string) (uid string, err error)
}
