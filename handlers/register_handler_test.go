package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/oleksiivelychko/go-account/initdb"
	"github.com/oleksiivelychko/go-account/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	initdb.LoadEnv()
	db, _ := initdb.TestDB()

	statement := "TRUNCATE accounts RESTART IDENTITY CASCADE"
	sqlExec := db.Exec(statement)
	if sqlExec.Error != nil {
		t.Errorf("[sql exec `"+statement+"`] -> %s", sqlExec.Error)
	}

	inputAccount := &models.Account{
		Email:    "test@test.test",
		Password: "secret",
	}
	payload := new(bytes.Buffer)
	_ = json.NewEncoder(payload).Encode(inputAccount)

	request, _ := http.NewRequest("POST", "/api/account/register", payload)
	response := httptest.NewRecorder()

	RegisterHandler(db)(response, request)

	responseBody := string(response.Body.Bytes())
	if responseBody == "email address already exists" {
		t.Fatalf(responseBody)
	}

	if response.Code != 201 {
		t.Fatalf("non-expected status code %v:\n\tbody: %v", "201", response.Code)
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
