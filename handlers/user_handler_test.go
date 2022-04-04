package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/requests"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler(t *testing.T) {
	db, _ := initTest()

	accountRepository := models.AccountRepository{DB: db, Debug: false}
	createdAccount, _ := accountRepository.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
	})

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

	userHandler := NewUserHandler(db)
	userHandler.ServeHTTP(response, request)

	if response.Code != 200 {
		t.Fatalf("non-expected status code %d", response.Code)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("unable to read response body: %s", err.Error())
	}

	account := &models.Account{}
	err = json.Unmarshal(body, &account)
	if err != nil {
		t.Fatalf("unable to unmarshal response body: %s", err.Error())
	}

	if account.Email != account.Email {
		t.Fatalf("emails doesn's match")
	}

	verifiedPassword := account.VerifyPassword("secret")
	if verifiedPassword != nil {
		t.Errorf("[func (model *Account) VerifyPassword(password string) error] -> %s", verifiedPassword)
	}
}
