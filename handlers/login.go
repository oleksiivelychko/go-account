package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/requests"
	"github.com/oleksiivelychko/go-account/services"
	"net/http"
)

type Login struct {
	accountService *services.Account
}

func NewLogin(accountService *services.Account) *Login {
	return &Login{accountService}
}

func (handler *Login) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	resp.Header().Set("Content-Type", "application/json")

	var inputAccount models.Account

	err := json.NewDecoder(req.Body).Decode(&inputAccount)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		_, _ = resp.Write([]byte(err.Error()))
		return
	}

	account, err := handler.accountService.Auth(inputAccount.Email, inputAccount.Password)
	if err != nil {
		resp.WriteHeader(http.StatusUnauthorized)
		_, _ = resp.Write([]byte(err.Error()))
		return
	}

	account, err = requests.AccessToken(account)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		_, _ = resp.Write([]byte(err.Error()))
		return
	}

	resp.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(resp).Encode(account)
}
