package infra

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/openhacku-saboten/OmnisCode-backend/config"
)

type Firebase struct {
	app *firebase.App
}

func NewFirebase() *Firebase {
	// Configが入っていればerrになりえないので無視
	app, _ := firebase.NewApp(
		context.Background(),
		&firebase.Config{
			DatabaseURL:      config.Firebase()["DatabaseURL"],
			ProjectID:        config.Firebase()["ProjectID"],
			ServiceAccountID: config.Firebase()["ServiceAccountID"],
			StorageBucket:    config.Firebase()["StorageBucket"],
		},
	)
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
