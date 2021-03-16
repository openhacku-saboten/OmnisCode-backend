package log

import "github.com/labstack/gommon/log"

// New はecho用のロガーを生成します
func New() *log.Logger {
	logger := log.New("api")
	logger.SetLevel(log.INFO)
	return logger
}
