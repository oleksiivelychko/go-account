package requests

import (
	"github.com/oleksiivelychko/go-account/initdb"
	"github.com/oleksiivelychko/go-account/models"
	"testing"
)

func TestRefreshTokenRequest(t *testing.T) {
	initdb.LoadEnv()

	account, err := AccessTokenRequest(&models.AccountSerialized{ID: 1})
	if err != nil {
		t.Fatalf("access token request error: %s", err)
	}

	_, err = RefreshTokenRequest(account)
	if err != nil {
		t.Fatalf("refresh token request error: %s", err)
	}
}
