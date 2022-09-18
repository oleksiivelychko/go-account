package requests

import (
	"encoding/json"
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"net/http"
	"os"
)

func AccessTokenRequest(accountSerialized *models.AccountSerialized) (*models.AccountSerialized, error) {
	apiRequestUrl := fmt.Sprintf("%s/access-token?userId=%d", os.Getenv("APP_JWT_URL"), accountSerialized.ID)
	response, err := http.Get(apiRequestUrl)
	if err != nil {
		return accountSerialized, fmt.Errorf("unable to make request to '%s'.\n%s", apiRequestUrl, err.Error())
	}

	if response.StatusCode != 201 {
		return accountSerialized, fmt.Errorf("unable to make successful request: status code is %d", response.StatusCode)
	}

	err = json.NewDecoder(response.Body).Decode(&accountSerialized)
	if err != nil {
		return accountSerialized, fmt.Errorf("unable to parse response body")
	}

	return accountSerialized, nil
}
