package requests

import (
	"github.com/oleksiivelychko/go-account/models"
	"testing"
)

func TestRequests_AccessToken(t *testing.T) {
	account, err := AccessToken(&models.AccountSerialized{ID: 1})
	if err != nil {
		t.Fatal(err)
	}

	if account.AccessToken == "" {
		t.Error("got empty accessToken")
	}
	if account.RefreshToken == "" {
		t.Error("got empty refreshToken")
	}
	if account.ExpirationTime == "" {
		t.Error("got empty expirationTime")
	}
}
