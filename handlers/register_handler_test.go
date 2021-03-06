package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/oleksiivelychko/go-account/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	db, _ := initTest()

	inputAccount := &models.Account{
		Email:    "test@test.test",
		Password: "secret",
	}
	payload := new(bytes.Buffer)
	_ = json.NewEncoder(payload).Encode(inputAccount)

	request, _ := http.NewRequest("POST", "/api/account/register", payload)
	response := httptest.NewRecorder()

	registerHandler := NewRegisterHandler(db)
	registerHandler.ServeHTTP(response, request)

	responseBody := string(response.Body.Bytes())
	if response.Code != 201 {
		t.Fatalf("non-expected status code: %d\nbody: %v", response.Code, responseBody)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("unable to read response body: %s", err.Error())
	}

	newAccount := &models.Account{}
	err = json.Unmarshal(body, &newAccount)
	if err != nil {
		t.Fatalf("unable to unmarshal response body: %s", err.Error())
	}

	if newAccount.Email != inputAccount.Email {
		t.Fatalf("emails doesn's match")
	}

	if newAccount.VerifyPassword(inputAccount.Password) != nil {
		t.Fatalf("passwords doesn't match")
	}
}
