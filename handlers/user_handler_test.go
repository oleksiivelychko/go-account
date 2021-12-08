package handlers

import (
	"github.com/oleksiivelychko/go-account/initdb"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/requests"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler(t *testing.T) {
	initdb.LoadEnv()
	db, _ := initdb.TestDB()

	accountSerialized, err := requests.AccessTokenRequest(&models.AccountSerialized{ID: 1})
	if err != nil {
		t.Fatalf("access token request error: %s", err)
	}

	request, _ := http.NewRequest("POST", "/api/account/user", nil)
	request.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{accountSerialized.AccessToken},
		"Expires":       []string{accountSerialized.ExpirationTime},
	}

	response := httptest.NewRecorder()

	UserHandler(db)(response, request)

	responseBody := string(response.Body.Bytes())
	if responseBody != "" {
		t.Fatalf("response has body: %s", responseBody)
	}

	if response.Code != 200 {
		t.Fatalf("non-expected status code %v:\n\tbody: %v", "200", response.Code)
	}
}
