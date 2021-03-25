package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CreateJSON は構造体をJSONにした時のフォーマットを確認するための関数です
// 実際のレスポンスがAPI仕様と同じかを確認するために利用することを想定しています
func CreateJSON(fileName string, entityData interface{}) error {
	fp, err := os.OpenFile(filepath.Clean(fileName), os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file: %s: %w", fileName, err)

	}
	defer func() {
		err := fp.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to close: %s", err.Error())
		}
	}()

	return MapEntity(fp, entityData)
}

// MapEntity はio.WriterにentityDataを書き込みます
func MapEntity(out io.Writer, entityData interface{}) error {
	err := json.NewEncoder(out).Encode(entityData)
	if err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}
	return nil
}
