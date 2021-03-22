package infra_test

import (
	"os"
	"testing"

	"github.com/openhacku-saboten/OmnisCode-backend/infra"
	"github.com/openhacku-saboten/OmnisCode-backend/log"
)

func TestNewDB(t *testing.T) {
	dbMap, err := infra.NewDB()
	if err != nil {
		t.Errorf("failed NewDB: %s", err.Error())
	}

	// きちんとCloseすること
	defer func() {
		err := dbMap.Db.Close()
		if err != nil {
			t.Errorf("failed to close DB: %s", err.Error())
		}
	}()
}

func ExampleNewDB() {
	logger := log.New()

	dbMap, err := infra.NewDB()
	if err != nil {
		logger.Errorf("failed NewDB: %s", err.Error())
		os.Exit(1)
	}

	// きちんとCloseすること
	defer func() {
		err := dbMap.Db.Close()
		if err != nil {
			logger.Errorf("failed to close DB: %s", err.Error())
		}
	}()
}
