package requests

import (
	"github.com/oleksiivelychko/go-account/models"
	"testing"
)

func TestRequests_RefreshToken(t *testing.T) {
	account, err := AccessToken(&models.AccountSerialized{ID: 1})
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = RefreshToken(account)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
