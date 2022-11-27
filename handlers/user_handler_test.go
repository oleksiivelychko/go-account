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
	"reflect"
	"testing"
)

func TestUserHandler(t *testing.T) {
	db, err := initdb.TestPrepare()
	if err != nil {
		t.Errorf("initialization test environment error: %s", err)
	}

	repository := repositories.NewRepository(db, false)
	accountRepository := repositories.NewAccountRepository(repository)
	accountService := services.NewAccountService(accountRepository)

	inputAccount := &models.Account{
		Email:    "test@test.test",
		Password: "secret",
	}
	err = accountRepository.Create(inputAccount)
	if err != nil {
		t.Errorf("unable to create account model: %s", err)
	}

	accountSerialized, err := requests.AccessTokenRequest(&models.AccountSerialized{ID: inputAccount.ID})
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

	modelAccount := &models.Account{}
	err = json.Unmarshal(body, &modelAccount)
	if err != nil {
		t.Fatalf("unable to unmarshal response body: %s", err.Error())
	}

	if modelAccount.ID != inputAccount.ID {
		t.Fatalf("id mismatch")
	}

	if modelAccount.Email != inputAccount.Email {
		t.Fatalf("email mismatch")
	}

	if reflect.TypeOf(modelAccount.CreatedAt).Name() != "DateTime" || reflect.TypeOf(modelAccount.UpdatedAt).Name() != "DateTime" {
		t.Fatalf("DateTime type mismatch")
	}
}
