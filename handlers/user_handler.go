package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/requests"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func UserHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenHeader := r.Header.Get("Authorization")
		expirationTime := r.Header.Get("Expires")

		accountSerialized := &models.AccountSerialized{
			AccessToken:    tokenHeader,
			ExpirationTime: expirationTime,
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

		accountRepository := models.AccountRepository{DB: db, Debug: false}
		account, err := accountRepository.FindOneByID(uint(userID))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(account)
	}
}
