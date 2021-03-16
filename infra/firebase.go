package infra

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
)

type Firebase struct {
	app *firebase.App
}

func NewFirebase() *Firebase {
	return &Firebase{}
}

func (f *Firebase) Authenticate(token string) (uid string, err error) {
	client, err := f.app.Auth(context.Background())
	if err != nil {
		return "", fmt.Errorf("error getting Auth client: %w", err)
	}
	authToken, err := client.VerifyIDTokenAndCheckRevoked(context.Background(), token)
	if err != nil {
		return "", fmt.Errorf("error verifying ID token: %w", err)
	}
	uid = authToken.UID
	return
}