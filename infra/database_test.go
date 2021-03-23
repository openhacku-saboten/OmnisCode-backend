package infra

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/openhacku-saboten/OmnisCode-backend/log"
)

func TestNewDB(t *testing.T) {
	dbMap, err := NewDB()
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

	dbMap, err := NewDB()
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

// truncateTable は指定したテーブルをtruncateするヘルパ関数です
func truncateTable(t *testing.T, dbMap *gorp.DbMap, tableName string) {
	t.Helper()

	// databaseを初期化する
	if _, err := dbMap.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		t.Fatal(err)
	}
	// タイミングの問題でTruncateが失敗することがあるので成功するまで試みる
	for i := 0; i < 5; i++ {
		_, err := dbMap.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName))
		if err == nil {
			break
		}
		if i == 4 {
			t.Fatal(err)
		}
		time.Sleep(time.Second * 1)
	}
	if _, err := dbMap.Exec("SET FOREIGN_KEY_CHECKS = 1"); err != nil {
		t.Fatal(err)
	}
}
