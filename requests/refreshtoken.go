package requests

import (
	"encoding/json"
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"net/http"
	"os"
)

func RefreshToken(accountSerialized *models.AccountSerialized) (*models.AccountSerialized, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/refresh-token", os.Getenv("APP_JWT_URL")), nil)
	if err != nil {
		return accountSerialized, err
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{accountSerialized.RefreshToken},
		"Expires`":      []string{accountSerialized.ExpirationTime},
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return accountSerialized, err
	}

	err = json.NewDecoder(resp.Body).Decode(&accountSerialized)
	if err != nil {
		return accountSerialized, err
	}

	return accountSerialized, nil
}
