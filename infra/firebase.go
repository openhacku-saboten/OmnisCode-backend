package infra

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/openhacku-saboten/OmnisCode-backend/config"
	"google.golang.org/api/option"
)

// NewFirebase はFirebase Authorization に接続するための構造体*auth.Clientを返す
func NewFirebase() (*auth.Client, error) {
	opt := option.WithCredentialsFile(config.GoogleAppCredentials())
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %w", err)
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %w", err)
	}
	return client, nil
}
