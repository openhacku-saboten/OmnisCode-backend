package infra

import (
	"context"

	firebase "firebase.google.com/go"
	"github.com/openhacku-saboten/OmnisCode-backend/config"
	"github.com/openhacku-saboten/OmnisCode-backend/log"
	"google.golang.org/api/option"
)

func NewFirebase() *firebase.App {
	logger := log.New()
	opt := option.WithCredentialsFile(config.GoogleAppCredentials())
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Errorf("error initializing app: %v\n", err)
	}
	return app
}
