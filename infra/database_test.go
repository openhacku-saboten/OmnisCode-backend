package infra_test

import (
	"github.com/openhacku-saboten/OmnisCode-backend/infra"
	"github.com/openhacku-saboten/OmnisCode-backend/log"
)

// TODO: example_testを書く
func ExampleNewDB() {
	logger := log.New()

	db, err := infra.NewDB()
	if err != nil {
		logger.Fatal(err)
	}
}
