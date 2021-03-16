package infra

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/openhacku-saboten/OmnisCode-backend/config"
	"github.com/openhacku-saboten/OmnisCode-backend/log"
	"google.golang.org/api/option"
)

type Firebase struct {
	app *firebase.App
}

func NewFirebase() *Firebase {
	logger := log.New()
	opt := option.WithCredentialsFile(config.GoogleAppCredentials())
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Errorf("error initializing app: %v\n", err)
	}
	return &Firebase{app: app}
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
