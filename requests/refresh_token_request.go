package requests

import (
	"encoding/json"
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"net/http"
	"os"
)

func RefreshTokenRequest(accountSerialized *models.AccountSerialized) (*models.AccountSerialized, error) {
	client := http.Client{}
	apiRequestUrl := fmt.Sprintf("%s/refresh-token", os.Getenv("APP_JWT_URL"))
	request, err := http.NewRequest("POST", apiRequestUrl, nil)
	if err != nil {
		return accountSerialized, err
	}

	request.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{accountSerialized.RefreshToken},
		"Expires`":      []string{accountSerialized.ExpirationTime},
	}

	response, err := client.Do(request)
	if err != nil {
		return accountSerialized, err
	}

	err = json.NewDecoder(response.Body).Decode(&accountSerialized)
	if err != nil {
		return accountSerialized, fmt.Errorf("unable to parse response body")
	}

	return accountSerialized, nil
}
