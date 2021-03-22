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
func (a *AuthRepository) Authenticate(ctx context.Context, token string) (uid string, err error) {
	authToken, err := a.firebase.VerifyIDTokenAndCheckRevoked(ctx, token)
	if err != nil {
		return "", fmt.Errorf("error verifying ID token: %w", err)
	}
	uid = authToken.UID
	return
}

// SetIconURL はuserIDからIconURLを取得して返す
func (a *AuthRepository) GetIconURL(ctx context.Context, uid string) (iconURL string, err error) {
	user, err := a.firebase.GetUser(ctx, uid)
	if err != nil {
		return "", fmt.Errorf("error getting user %s from firebase: %w", uid, err)
	}
	iconURL = user.PhotoURL
	return
}
