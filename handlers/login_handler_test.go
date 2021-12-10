package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/oleksiivelychko/go-account/initdb"
	"github.com/oleksiivelychko/go-account/models"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initTest() (*gorm.DB, error) {
	initdb.LoadEnv()
	db, err := initdb.TestDB()
	err = models.AutoMigrate(db)

	statement := "TRUNCATE accounts RESTART IDENTITY CASCADE"
	sqlExec := db.Exec(statement)
	if sqlExec.Error != nil {
		return nil, sqlExec.Error
	}

	statement = "TRUNCATE roles RESTART IDENTITY CASCADE"
	sqlExec = db.Exec(statement)
	if sqlExec.Error != nil {
		return nil, sqlExec.Error
	}

	return db, err
}

func TestLoginHandler(t *testing.T) {
	db, _ := initTest()

	accountRepository := models.AccountRepository{DB: db, Debug: false}
	_, _ = accountRepository.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	inputAccount := &models.Account{
		Email:    "test@test.test",
		Password: "secret",
	}

	payload := new(bytes.Buffer)
	_ = json.NewEncoder(payload).Encode(inputAccount)

	request, _ := http.NewRequest("POST", "/api/account/login", payload)
	response := httptest.NewRecorder()

	LoginHandler(db)(response, request)

	responseBody := string(response.Body.Bytes())
	if responseBody == "invalid password" {
		t.Fatalf(responseBody)
	}

	if response.Code != 200 {
		t.Fatalf("non-expected status code %v:\n\tbody: %v", "200", response.Code)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("unable to read response body: %s", err.Error())
	}

	loggedAccount := &models.AccountSerialized{}
	err = json.Unmarshal(body, &loggedAccount)
	if err != nil {
		t.Fatalf("unable to unmarshal response body: %s", err.Error())
	}

	if loggedAccount.Email != inputAccount.Email {
		t.Fatalf("emails doesn's match")
	}

	if loggedAccount.AccessToken == "" {
		t.Fatalf("got empty `access-token`")
	}

	if loggedAccount.RefreshToken == "" {
		t.Fatalf("got empty `refresh-token`")
	}

	if loggedAccount.ExpirationTime == "" {
		t.Fatalf("got empty `expiration-time`")
	}
}
