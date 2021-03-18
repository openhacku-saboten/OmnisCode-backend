package infra

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/openhacku-saboten/OmnisCode-backend/config"
	"github.com/openhacku-saboten/OmnisCode-backend/log"

	_ "github.com/go-sql-driver/mysql"
)

// NewDB はMySQLサーバに接続して、*gorp.DbMapを生成します
func NewDB() (*gorp.DbMap, error) {
	db, err := sql.Open("mysql", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL: %w", err)
	}

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	logger := log.New()

	for {
		err := db.Ping()
		if err == nil {
			break
		}
		logger.Infof("%s\n", err.Error())
		time.Sleep(time.Second * 2)
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	logger.Info("DB Ready!")
	return dbMap, nil
}
