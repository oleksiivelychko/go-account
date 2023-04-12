package requests

import (
	"github.com/joho/godotenv"
	"github.com/oleksiivelychko/go-account/models"
	"log"
	"testing"
)

func TestRequests_RefreshToken(t *testing.T) {
	err := godotenv.Load("./../.env.test")
	if err != nil {
		log.Fatal("error loading .env.test file")
	}

	account, err := AccessToken(&models.AccountSerialized{ID: 1})
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = RefreshToken(account)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
