package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"net/http"
	"os"
)

func AuthorizeToken(accountSerialized *models.AccountSerialized) (*models.AccountSerialized, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/authorize-token", os.Getenv("APP_JWT_URL")), nil)
	if err != nil {
		return accountSerialized, err
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{accountSerialized.AccessToken},
		"Expires":       []string{accountSerialized.ExpirationTime},
	}

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode == http.StatusOK {
		return accountSerialized, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return RefreshToken(accountSerialized)
	}

	err = json.NewDecoder(resp.Body).Decode(&accountSerialized)
	if err != nil {
		return accountSerialized, err
	}

	return accountSerialized, errors.New(accountSerialized.ErrorMessage)
}
