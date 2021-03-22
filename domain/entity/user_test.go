package entity

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUser_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr error
	}{
		{
			name:    "問題なければnilを返す",
			user:    NewUser("id", "name", "profile", "twitter", "url"),
			wantErr: nil,
		},
		{
			name:    "IDが空ならエラー",
			user:    NewUser("", "name", "profile", "twitter", "url"),
			wantErr: errors.New("user ID must not be empty"),
		},
		{
			name:    "Nameが空ならエラー",
			user:    NewUser("id", "", "profile", "twitter", "url"),
			wantErr: ErrEmptyUserName,
		},
		{
			name:    "IDが129文字以上ならエラー",
			user:    NewUser(strings.Repeat("a", 129), "name", "profile", "twitter", "url"),
			wantErr: ErrTooLong,
		},
		{
			name:    "Nameが129文字以上ならエラー",
			user:    NewUser("id", strings.Repeat("a", 129), "profile", "twitter", "url"),
			wantErr: ErrTooLong,
		},
		{
			name:    "TwitterIDが16文字以上ならエラー",
			user:    NewUser("id", "name", "profile", strings.Repeat("a", 16), "url"),
			wantErr: ErrTooLong,
		},
		{
			name:    "マルチバイト文字列は１文字とカウント",
			user:    NewUser("id", strings.Repeat("あ", 128), "profile", "twitter", "url"),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.IsValid()
			if got == nil && tt.wantErr == nil {
				return
			}
			if got.Error() != tt.wantErr.Error() {
				t.Errorf("User.IsValid() = %s, want = %s", got.Error(), tt.wantErr.Error())
			}
		})
	}
}

func TestUser_Format(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		wantUser *User
	}{
		{
			name:     "TwitterIDが空なら何もしない",
			user:     NewUser("id", "name", "profile", "", "url"),
			wantUser: NewUser("id", "name", "profile", "", "url"),
		},
		{
			name:     "TwitterIDが@が含まれていないなら何もしない",
			user:     NewUser("id", "name", "profile", "twitter", "url"),
			wantUser: NewUser("id", "name", "profile", "twitter", "url"),
		},
		{
			name:     "TwitterIDに@が含まれているなら取り除く",
			user:     NewUser("id", "name", "profile", "@twitter", "url"),
			wantUser: NewUser("id", "name", "profile", "twitter", "url"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.user.Format()

			if diff := cmp.Diff(tt.wantUser, tt.user); diff != "" {
				t.Errorf("Data (-want +got) =\n%s\n", diff)
			}
		})
	}
}
