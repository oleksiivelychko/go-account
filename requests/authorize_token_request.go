package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"net/http"
	"os"
)

func AuthorizeTokenRequest(accountSerialized *models.AccountSerialized) (*models.AccountSerialized, error) {
	client := http.Client{}
	apiRequestUrl := fmt.Sprintf("%s/authorize-token", os.Getenv("APP_JWT_URL"))
	request, err := http.NewRequest("POST", apiRequestUrl, nil)
	if err != nil {
		return accountSerialized, err
	}

	request.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{accountSerialized.AccessToken},
		"Expires":       []string{accountSerialized.ExpirationTime},
	}

	response, err := client.Do(request)
	if err != nil {
		return accountSerialized, err
	}

	if response.StatusCode == 200 {
		return accountSerialized, nil
	}

	err = json.NewDecoder(response.Body).Decode(&accountSerialized)
	if err != nil {
		return accountSerialized, fmt.Errorf("unable to parse response body")
	}

	if accountSerialized.ErrorCode == issuer.FailedExpirationTimeClaim {
		return RefreshTokenRequest(accountSerialized)
	}

	return accountSerialized, errors.New(accountSerialized.ErrorMessage)
}
