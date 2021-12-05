package requests

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func AccessTokenRequest(userId uint) (*http.Response, error) {
	apiAccessTokenUrl := os.Getenv("API_ACCESS_TOKEN_URL")
	if apiAccessTokenUrl == "" {
		return nil, errors.New("API_ACCESS_TOKEN_URL is not set")
	}

	var apiRequestUrl = fmt.Sprintf("%s?userId=%d", apiAccessTokenUrl, userId)
	response, err := http.Get(apiRequestUrl)
	if err != nil {
		return response, errors.New(fmt.Sprintf("unable to make request to `%s`", apiRequestUrl))
	}

	return response, err
}
