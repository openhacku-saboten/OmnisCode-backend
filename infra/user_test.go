package infra

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
)

func TestUserRepository_FindByID(t *testing.T) {
	dbMap, err := NewDB()
	if err != nil {
		t.Fatalf(err.Error())
	}
	dbMap.AddTableWithName(UserDTO{}, "users")
	truncateTable(t, dbMap, "users")
	if err := dbMap.Insert(&UserDTO{
		ID:        "existing-id",
		Name:      "existingUser",
		Profile:   "existing",
		TwitterID: "existing",
	}); err != nil {
		t.Fatal(err)
	}

	userRepo := NewUserRepository(dbMap)

	tests := []struct {
		name     string
		userID   string
		wantUser *entity.User
		wantErr  error
	}{
		{
			name:     "正しくユーザーを取得できる",
			userID:   "existing-id",
			wantUser: entity.NewUser("existing-id", "existingUser", "existing", "existing", ""),
			wantErr:  nil,
		},
		{
			name:     "存在しないユーザーの場合はErrNoRows",
			userID:   "not-existing-id",
			wantUser: nil,
			wantErr:  entity.ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := userRepo.FindByID(tt.userID)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				diff := cmp.Diff(tt.wantUser, gotUser)
				if diff != "" {
					t.Errorf("Data (-want +got) =\n%s\n", diff)
				}
			}
		})
	}
}

func TestUserRepository_Insert(t *testing.T) {
	dbMap, err := NewDB()
	if err != nil {
		t.Fatalf(err.Error())
	}
	dbMap.AddTableWithName(UserDTO{}, "users")
	truncateTable(t, dbMap, "users")
	if err := dbMap.Insert(&UserDTO{
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
