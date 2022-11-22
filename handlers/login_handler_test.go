package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/oleksiivelychko/go-account/initdb"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/repositories"
	"github.com/oleksiivelychko/go-account/services"
	"gorm.io/gorm"
	"io"
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

	inputAccount := &models.Account{
		Email:    "test@test.test",
		Password: "secret",
	}

	accountRepository := repositories.NewAccountRepository(db, false)
	accountService := services.NewAccountService(accountRepository)

	inputAccount, err := accountService.Create(&models.Account{
		Email:    inputAccount.Email,
		Password: inputAccount.Password,
	})
	if err != nil {
		t.Fatalf("unable to create account: %s", err)
	}

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

	loggedAccount := &models.AccountSerialized{}
	err = json.Unmarshal(body, &loggedAccount)
	if err != nil {
		t.Fatalf("unable to unmarshal response body: %s", err.Error())
	}

	if loggedAccount.Email != inputAccount.Email {
		t.Fatalf("email mismatch")
	}

	if loggedAccount.AccessToken == "" {
		t.Fatalf("got empty `access_token`")
	}

	if loggedAccount.RefreshToken == "" {
		t.Fatalf("got empty `refresh_token`")
	}

	if loggedAccount.ExpirationTime == "" {
		t.Fatalf("got empty `expiration_time`")
	}
}
