package config_test

import (
	"testing"

	"github.com/openhacku-saboten/OmnisCode-backend/config"
)

func TestPort(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		want string
	}{
		{
			name: "正しくポートを取得できる",
			want: "8080",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := config.Port(); got != tc.want {
				t.Errorf("Port() = %s, want = %s", got, tc.want)
			}
		})
	}
}

func TestDSN(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		want string
	}{
		{
			name: "正しくDSNを取得できる",
			want: "test:test@tcp(localhost:3306)/test?parseTime=true&collation=utf8mb4_bin",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := config.DSN(); got != tc.want {
				t.Errorf("DSN() = %s, want = %s", got, tc.want)
			}
		})
	}
}

func TestGoogleAppCredentials(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		want string
	}{
		{
			name: "正しくGoogleAppCredentialsを取得できる",
			want: "firebaseCredentials.json",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := config.GoogleAppCredentials(); got != tc.want {
				t.Errorf("GoogleAppCredentials() = %s, want = %s", got, tc.want)
			}
		})
	}
}
