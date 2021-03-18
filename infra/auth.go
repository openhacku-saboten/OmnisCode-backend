package infra

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
)

type AuthRepository struct {
	fb *firebase.App
}

func NewAuthRepository(fb *firebase.App) *AuthRepository {
	return &AuthRepository{fb: fb}
}

func (a *AuthRepository) Authenticate(token string) (uid string, err error) {
	client, err := a.fb.Auth(context.Background())
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
