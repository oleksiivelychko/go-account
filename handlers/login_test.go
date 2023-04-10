package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/oleksiivelychko/go-account/db"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/repositories"
	"github.com/oleksiivelychko/go-account/services"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers_Login(t *testing.T) {
	sessionDB, err := db.PrepareTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := repositories.NewRepository(sessionDB, false)
	accountRepo := repositories.NewAccount(repo)
	accountService := services.NewAccount(accountRepo)

	inputAccount, err := accountService.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	if err != nil {
		t.Fatal(err)
	}

	// send non-hashed password
	inputAccount.Password = "secret"

	payload := new(bytes.Buffer)
	_ = json.NewEncoder(payload).Encode(inputAccount)

	req, _ := http.NewRequest("POST", "/api/account/login", payload)
	resp := httptest.NewRecorder()

	loginHandler := NewLogin(accountService)
	loginHandler.ServeHTTP(resp, req)

	respBody := string(resp.Body.Bytes())
	if resp.Code != http.StatusOK {
		t.Fatalf("non-expected status code: %d\nbody: %v", resp.Code, respBody)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf(err.Error())
	}

	accountSerialized := &models.AccountSerialized{}
	err = json.Unmarshal(body, &accountSerialized)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if accountSerialized.ID != inputAccount.ID {
		t.Error("id mismatch")
	}

	if accountSerialized.Email != inputAccount.Email {
		t.Error("email mismatch")
	}

	if accountSerialized.AccessToken == "" {
		t.Error("got empty accessToken")
	}

	if accountSerialized.RefreshToken == "" {
		t.Error("got empty refreshToken")
	}

	if accountSerialized.ExpirationTime == "" {
		t.Error("got empty expirationTime")
	}
}
