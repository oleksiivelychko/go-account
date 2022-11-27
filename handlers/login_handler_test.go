package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/oleksiivelychko/go-account/initdb"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/repositories"
	"github.com/oleksiivelychko/go-account/services"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	db, err := initdb.TestPrepare()
	if err != nil {
		t.Errorf("initialization test environment error: %s", err)
	}

	repository := repositories.NewRepository(db, false)
	accountRepository := repositories.NewAccountRepository(repository)
	accountService := services.NewAccountService(accountRepository)

	inputAccount, err := accountService.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("unable to create account model: %s", err)
	}

	// send non-hashed password
	inputAccount.Password = "secret"

	payload := new(bytes.Buffer)
	_ = json.NewEncoder(payload).Encode(inputAccount)

	request, _ := http.NewRequest("POST", "/api/account/login", payload)
	response := httptest.NewRecorder()

	loginHandler := NewLoginHandler(accountService)
	loginHandler.ServeHTTP(response, request)

	responseBody := string(response.Body.Bytes())
	if response.Code != 200 {
		t.Fatalf("non-expected status code: %d\nbody: %v", response.Code, responseBody)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("unable to read response body: %s", err.Error())
	}

	accountSerialized := &models.AccountSerialized{}
	err = json.Unmarshal(body, &accountSerialized)
	if err != nil {
		t.Fatalf("unable to unmarshal response body: %s", err.Error())
	}

	if accountSerialized.ID != inputAccount.ID {
		t.Fatalf("id mismatch")
	}

	if accountSerialized.Email != inputAccount.Email {
		t.Fatalf("email mismatch")
	}

	if accountSerialized.AccessToken == "" {
		t.Fatalf("got empty `access_token`")
	}

	if accountSerialized.RefreshToken == "" {
		t.Fatalf("got empty `refresh_token`")
	}

	if accountSerialized.ExpirationTime == "" {
		t.Fatalf("got empty `expiration_time`")
	}
}
