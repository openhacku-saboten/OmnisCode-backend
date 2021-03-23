package infra

import (
	"errors"
	"testing"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

func TestPostInfra_Insert(t *testing.T) {
	dbMap, err := NewDB()
	if err != nil {
		t.Fatalf(err.Error())
	}
	dbMap.AddTableWithName(PostDTO{}, "posts")
	truncateUser(t, dbMap)
	// TODO:
	if err := dbMap.Insert(&PostDTO{
		ID:        "existing-id",
		Name:      "existingUser",
		Profile:   "existing",
		TwitterID: "existing",
	}); err != nil {
		t.Fatal(err)
	}

	userRepo := NewUserRepository(dbMap)

	tests := []struct {
		name    string
		user    *entity.User
		wantErr error
	}{
		{
			name:    "すでに存在するユーザーIDならErrDuplicatedUser",
			user:    entity.NewUser("existing-id", "newUser", "new", "new", ""),
			wantErr: entity.ErrDuplicatedUser,
		},
		{
			name:    "すでに存在するTwitterIDならErrDuplicatedTwitterID",
			user:    entity.NewUser("new-id", "newUser", "new", "existing", ""),
			wantErr: entity.ErrDuplicatedTwitterID,
		},
		{
			name:    "正しくユーザーを作成できる",
			user:    entity.NewUser("new-id", "newUser", "new", "new", ""),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := userRepo.Insert(tt.user)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
		})
	}
}

func truncateUser(t *testing.T, dbMap *gorp.DbMap) {
	t.Helper()

	// databaseを初期化する
	if _, err := dbMap.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		t.Fatal(err)
	}
	// タイミングの問題でTruncateが失敗することがあるので成功するまで試みる
	for i := 0; i < 5; i++ {
		_, err := dbMap.Exec("TRUNCATE TABLE users")
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
