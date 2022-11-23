package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/oleksiivelychko/go-account/initdb"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/repositories"
	"github.com/oleksiivelychko/go-account/requests"
	"github.com/oleksiivelychko/go-account/services"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler(t *testing.T) {
	db, err := initdb.TestPrepare()
	if err != nil {
		t.Errorf("initialization test environment error: %s", err)
	}

	accountRepository := repositories.NewAccountRepository(db, false)
	createdAccount, _ := accountRepository.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	accountService := services.NewAccountService(accountRepository)

	accountSerialized, err := requests.AccessTokenRequest(&models.AccountSerialized{ID: createdAccount.ID})
	if err != nil {
		t.Fatalf("access token request error: %s", err)
	}

	requestUrl := fmt.Sprintf("/api/account/user/?userId=%d", accountSerialized.ID)
	request, _ := http.NewRequest("POST", requestUrl, nil)
	request.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{accountSerialized.AccessToken},
		"Expires":       []string{accountSerialized.ExpirationTime},
	}

	response := httptest.NewRecorder()

	userHandler := NewUserHandler(accountService)
	userHandler.ServeHTTP(response, request)

	if response.Code != 200 {
		t.Fatalf("non-expected status code %d", response.Code)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("unable to read response body: %s", err.Error())
	}

	account := &models.Account{}
	err = json.Unmarshal(body, &account)
	if err != nil {
		t.Fatalf("unable to unmarshal response body: %s", err.Error())
	}

	if account.Email != account.Email {
		t.Fatalf("email mismatch")
	}

	verifiedPasswordErr := accountService.VerifyPassword(account, "secret")
	if verifiedPasswordErr != nil {
		t.Errorf("unable to verify password: %s", verifiedPasswordErr)
	}
}
