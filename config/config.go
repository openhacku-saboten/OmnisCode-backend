package config

import (
	"fmt"
	"os"
)

// Port は環境変数に書かれているPORTの値をstringで返す関数です
func Port() string {
	return os.Getenv("PORT")
}

// DSN は環境変数の情報をもとに、DataSourceNameを返す関数です
func DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	) + "?parseTime=true&collation=utf8mb4_bin"
}

func Firebase() map[string]string {
	return map[string]string{
		"DatabaseURL":      os.Getenv("FIREBASE_URL"),
		"ProjectID":        os.Getenv("FIREBASE_PROJECT_ID"),
		"ServiceAccountID": os.Getenv("FIREBASE_SB_ID"),
		"StorageBucket":    os.Getenv("FIREBASE_STORAGE_BUCKET"),
	}
}
