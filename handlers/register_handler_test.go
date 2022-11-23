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

func TestRegisterHandler(t *testing.T) {
	db, err := initdb.TestPrepare()
	if err != nil {
		t.Errorf("initialization test environment error: %s", err)
	}

	inputAccount := &models.Account{
		Email:    "test@test.test",
		Password: "secret",
	}
	payload := new(bytes.Buffer)
	_ = json.NewEncoder(payload).Encode(inputAccount)

	request, _ := http.NewRequest("POST", "/api/account/register", payload)
	response := httptest.NewRecorder()

	accountRepository := repositories.NewAccountRepository(db, false)
	roleRepository := repositories.NewRoleRepository(db, false)
	accountService := services.NewAccountService(accountRepository)
	roleService := services.NewRoleService(roleRepository)

	registerHandler := NewRegisterHandler(accountService, roleService)
	registerHandler.ServeHTTP(response, request)

	responseBody := string(response.Body.Bytes())
	if response.Code != 201 {
		t.Fatalf("non-expected status code: %d\nbody: %v", response.Code, responseBody)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("unable to read response body: %s", err.Error())
	}

	newAccount := &models.Account{}
	err = json.Unmarshal(body, &newAccount)
	if err != nil {
		t.Fatalf("unable to unmarshal response body: %s", err.Error())
	}

	if newAccount.Email != inputAccount.Email {
		t.Fatalf("email mismatch")
	}

	if accountService.VerifyPassword(inputAccount, inputAccount.Password) != nil {
		t.Fatalf("password mismatch")
	}
}
