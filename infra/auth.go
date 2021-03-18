package infra

import (
	"context"
	"fmt"

	"firebase.google.com/go/auth"
)

type AuthRepository struct {
	firebase *auth.Client
}

func NewAuthRepository(firebase *auth.Client) *AuthRepository {
	return &AuthRepository{firebase: firebase}
}

// Authenticate はTokenをfirebaseに照合してuserIDを返す
func (a *AuthRepository) Authenticate(token string) (uid string, err error) {
	authToken, err := a.firebase.VerifyIDTokenAndCheckRevoked(context.Background(), token)
	if err != nil {
		return "", fmt.Errorf("error verifying ID token: %w", err)
	}
	uid = authToken.UID
	return
}
