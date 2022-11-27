package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/requests"
	"github.com/oleksiivelychko/go-account/services"
	"net/http"
)

type LoginHandler struct {
	accountService *services.AccountService
}

func NewLoginHandler(s *services.AccountService) *LoginHandler {
	return &LoginHandler{s}
}

// type Handler interface
func (handler *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var inputAccount models.Account

	err := json.NewDecoder(r.Body).Decode(&inputAccount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	account, err := handler.accountService.Auth(inputAccount.Email, inputAccount.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	account, err = requests.AccessTokenRequest(account)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(account)
}
