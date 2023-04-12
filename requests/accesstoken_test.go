package requests

import (
	"github.com/joho/godotenv"
	"github.com/oleksiivelychko/go-account/models"
	"log"
	"testing"
)

func TestRequests_AccessToken(t *testing.T) {
	err := godotenv.Load("./../.env.test")
	if err != nil {
		log.Fatal("error loading .env.test file")
	}

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
