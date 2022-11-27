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

type UserHandler struct {
	accountService *services.AccountService
}

func NewUserHandler(s *services.AccountService) *UserHandler {
	return &UserHandler{s}
}

func (handler *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	accountSerialized := &models.AccountSerialized{
		AccessToken:    r.Header.Get("Authorization"),
		ExpirationTime: r.Header.Get("Expires"),
	}

	accountSerialized, err := requests.AuthorizeTokenRequest(accountSerialized)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(accountSerialized)
		return
	}

	v := r.URL.Query()
	userID, err := strconv.ParseInt(v.Get("userId"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "unable to get user identifier as `userId` from URL query: %s", err.Error())
		return
	}

	modelAccount, err := handler.accountService.GetRepository().FindOneByID(uint(userID))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(modelAccount)
}
