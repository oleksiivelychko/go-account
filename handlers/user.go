package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/requests"
	"github.com/oleksiivelychko/go-account/services"
	"net/http"
	"strconv"
)

type User struct {
	accountService *services.Account
}

func NewUser(accountService *services.Account) *User {
	return &User{accountService}
}

func (handler *User) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	accountSerialized := &models.AccountSerialized{
		AccessToken:    req.Header.Get("Authorization"),
		ExpirationTime: req.Header.Get("Expires"),
	}

	accountSerialized, err := requests.AuthorizeToken(accountSerialized)
	if err != nil {
		resp.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(resp).Encode(accountSerialized)
		return
	}

	queryValues := req.URL.Query()
	userID, err := strconv.ParseInt(queryValues.Get("userID"), 10, 64)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(resp, "unable to get userID from URL: %s", err.Error())
		return
	}

	modelAccount, err := handler.accountService.GetRepository().FindOneByID(uint(userID))
	if err != nil {
		resp.WriteHeader(http.StatusUnauthorized)
		_, _ = resp.Write([]byte(err.Error()))
		return
	}

	resp.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(resp).Encode(modelAccount)
}
