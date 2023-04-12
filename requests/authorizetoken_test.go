package requests

import (
	"github.com/joho/godotenv"
	"github.com/oleksiivelychko/go-account/models"
	"log"
	"testing"
	"time"
)

func TestRequests_AuthorizeToken(t *testing.T) {
	err := godotenv.Load("./../.env.test")
	if err != nil {
		log.Fatal("error loading .env.test file")
	}

	account, err := AccessToken(&models.AccountSerialized{ID: 1})
	if err != nil {
		t.Fatal(err)
	}

	_, err = AuthorizeToken(account)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRequest_AuthorizeByExpiredToken(t *testing.T) {
	account, err := AccessToken(&models.AccountSerialized{ID: 1})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	_, err = AuthorizeToken(account)
	if err != nil {
		t.Fatal(err)
	}
}
