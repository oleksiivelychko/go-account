package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/requests"
	"gorm.io/gorm"
	"net/http"
)

type LoginHandler struct {
	db *gorm.DB
}

func NewLoginHandler(db *gorm.DB) *LoginHandler {
	return &LoginHandler{db}
}

// type Handler interface
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var inputAccount models.Account
	err := json.NewDecoder(r.Body).Decode(&inputAccount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accountRepository := models.AccountRepository{DB: h.db, Debug: false}
	account, err := accountRepository.MakeAuth(inputAccount.Email, inputAccount.Password)

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
