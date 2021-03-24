package service

import (
	"encoding/json"
	"fmt"
	"os"
)

// CreateJSONは構造体をJSONにした時のフォーマットを確認するための関数です
// 実際のレスポンスがAPI仕様と同じかを確認するために利用することを想定しています
func CreateJSON(fileName string, entityData interface{}) error {
	fp, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file: %s: %w", fileName, err)

	}
	defer fp.Close()

	err = json.NewEncoder(fp).Encode(entityData)
	if err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}
	return nil
}
