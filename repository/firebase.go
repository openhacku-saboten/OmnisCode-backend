package repository

type Firebase interface {
	Authenticate(token string) (uid string, err error)
}