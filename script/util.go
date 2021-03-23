package script

import (
	"encoding/json"
	"os"
	"time"

	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

func CreateJson() {
	fp, _ := os.OpenFile("test.json", os.O_CREATE|os.O_WRONLY, 0600)
	defer fp.Close()

	ve := &entity.Post{
		ID:        0,
		UserID:    "testID",
		Title:     "test title",
		Code:      "package main\n\nimport \"fmt\"\n\nfunc main(){fmt.Println(\"This is test.\")}",
		Language:  "Go",
		Content:   "Test code",
		Source:    "github.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	json.NewEncoder(fp).Encode(ve)
}
