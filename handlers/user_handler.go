package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/requests"
	"gorm.io/gorm"
	"net/http"
)

func UserHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenHeader := r.Header.Get("Authorization")
		expirationTime := r.Header.Get("Expires")

		account := &models.AccountSerialized{
			AccessToken:    tokenHeader,
			ExpirationTime: expirationTime,
		}

		account, err := requests.AuthorizeTokenRequest(account)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(account)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
