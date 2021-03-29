package infra

import (
	"context"
	"fmt"

	"firebase.google.com/go/auth"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

// AuthRepository は認証情報の永続化と再構成のためのリポジトリです
type AuthRepository struct {
	firebase *auth.Client
}

// NewAuthRepository は認証情報のリポジトリのポインタを生成する関数です
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

// GetIconURL はuserIDからIconURLを取得して返す
func (a *AuthRepository) GetIconURL(ctx context.Context, uid string) (iconURL string, err error) {
	user, err := a.firebase.GetUser(ctx, uid)
	if err != nil {
		return "", fmt.Errorf("error getting user %s from firebase: %w", uid, err)
	}
	iconURL = user.PhotoURL
	return
}

// Delete はuserIDからuserを削除します
func (a *AuthRepository) Delete(ctx context.Context, user *entity.User) error {
	if err := a.firebase.DeleteUser(ctx, user.ID); err != nil {
		return fmt.Errorf("failed firebase.DeleteUser: %w", err)
	}
	return nil
}
