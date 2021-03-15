package infra

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/task4233/techtrain-mission/gameapi/config"
)

// NewDB はMySQLサーバに接続して、*gorp.DbMapを生成します
func NewDB() (*gorp.DbMap, error) {
	db, err := sql.Open("mysql", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL: %w", err)
	}

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	for {
		err := db.Ping()
		if err == nil {
			break
		}
		time.Sleep(time.Second * 2)
	}
}
