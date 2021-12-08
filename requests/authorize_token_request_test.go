package requests

import (
	"github.com/oleksiivelychko/go-account/initdb"
	"github.com/oleksiivelychko/go-account/models"
	"testing"
	"time"
)

func TestAuthorizeTokenRequest(t *testing.T) {
	initdb.LoadEnv()

	account, err := AccessTokenRequest(&models.AccountSerialized{ID: 1})
	if err != nil {
		t.Fatalf("access token request error: %s", err)
	}

	_, err = AuthorizeTokenRequest(account)
	if err != nil {
		t.Fatalf("authorize token request error: %s", err)
	}
}

func TestExpiredAuthorizeTokenRequest(t *testing.T) {
	initdb.LoadEnv()

	account, err := AccessTokenRequest(&models.AccountSerialized{ID: 1})
	if err != nil {
		t.Fatalf("access token request error: %s", err)
	}

	time.Sleep(61 * time.Second)

	_, err = AuthorizeTokenRequest(account)
	if err != nil {
		t.Fatalf("authorize token request error: %s", err)
	}
}
