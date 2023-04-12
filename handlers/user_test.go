package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/oleksiivelychko/go-account/db"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/repositories"
	"github.com/oleksiivelychko/go-account/requests"
	"github.com/oleksiivelychko/go-account/services"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandler_User(t *testing.T) {
	sessionDB, err := db.PrepareTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := repositories.NewRepository(sessionDB, false)
	accountRepo := repositories.NewAccount(repo)
	accountService := services.NewAccount(accountRepo)

	inputAccount := &models.Account{
		Email:    "test@test.test",
		Password: "secret",
	}

	err = accountRepo.Create(inputAccount)
	if err != nil {
		t.Fatal(err)
	}

	accountSerialized, err := requests.AccessToken(&models.AccountSerialized{ID: inputAccount.ID})
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("/api/account/user/?userID=%d", accountSerialized.ID)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{accountSerialized.AccessToken},
		"Expires":       []string{accountSerialized.ExpirationTime},
	}

	resp := httptest.NewRecorder()

	userHandler := NewUser(accountService)
	userHandler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("non-expected status code %d", resp.Code)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	modelAccount := &models.Account{}
	err = json.Unmarshal(body, &modelAccount)
	if err != nil {
		t.Fatal(err.Error())
	}

	if modelAccount.ID != inputAccount.ID {
		t.Error("id mismatch")
	}

	if modelAccount.Email != inputAccount.Email {
		t.Error("email mismatch")
	}

	if reflect.TypeOf(modelAccount.CreatedAt).Name() != "DateTime" || reflect.TypeOf(modelAccount.UpdatedAt).Name() != "DateTime" {
		t.Error("datetime type mismatch")
	}
}
