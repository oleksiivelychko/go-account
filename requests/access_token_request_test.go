package requests

import (
	"github.com/oleksiivelychko/go-account/initdb"
	"github.com/oleksiivelychko/go-account/models"
	"testing"
)

func TestAccessTokenRequest(t *testing.T) {
	initdb.LoadEnv()

	account, err := AccessTokenRequest(&models.AccountSerialized{ID: 1})
	if err != nil {
		t.Fatalf("access token request error: %s", err)
	}

	if account.AccessToken == "" {
		t.Fatalf("got empty `access_token`")
	}
	if account.RefreshToken == "" {
		t.Fatalf("got empty `refresh_token`")
	}
	if account.ExpirationTime == "" {
		t.Fatalf("got empty `expiration_time`")
	}
}
