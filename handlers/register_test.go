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

func TestRegisterHandler(t *testing.T) {
	sessionDB, err := db.PrepareTestDB()
	if err != nil {
		t.Fatal(err)
	}

	inputAccount := &models.Account{
		Email:    "test@test.test",
		Password: "secret",
	}

	payload := new(bytes.Buffer)
	_ = json.NewEncoder(payload).Encode(inputAccount)

	req, _ := http.NewRequest("POST", "/api/account/register", payload)
	resp := httptest.NewRecorder()

	repo := repositories.NewRepository(sessionDB, false)
	accountRepo := repositories.NewAccount(repo)
	roleRepo := repositories.NewRole(repo)
	accountService := services.NewAccount(accountRepo)
	roleService := services.NewRole(roleRepo)

	registerHandler := NewRegister(accountService, roleService)
	registerHandler.ServeHTTP(resp, req)

	respBody := string(resp.Body.Bytes())
	if resp.Code != 201 {
		t.Fatalf("non-expected status code: %d\nbody: %v", resp.Code, respBody)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	accountSerialized := &models.AccountSerialized{}
	err = json.Unmarshal(body, &accountSerialized)
	if err != nil {
		t.Fatal(err.Error())
	}

	if accountSerialized.Email != inputAccount.Email {
		t.Fatal("email mismatch")
	}
}
