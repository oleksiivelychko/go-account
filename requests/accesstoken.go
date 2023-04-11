package requests

import (
	"encoding/json"
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"net/http"
	"os"
)

func AccessToken(account *models.AccountSerialized) (*models.AccountSerialized, error) {
	resp, err := http.Get(fmt.Sprintf("%s/access-token?userID=%d", os.Getenv("APP_JWT_URL"), account.ID))
	if err != nil {
		return account, err
	}

	if resp.StatusCode != 201 {
		return account, fmt.Errorf("non-expected status code %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return account, err
	}

	return account, nil
}
