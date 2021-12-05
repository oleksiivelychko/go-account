package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/requests"
	"gorm.io/gorm"
	"net/http"
)

func LoginHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		accountRepository := models.AccountRepository{DB: db, Debug: true}
		account, err := accountRepository.MakeAuth(inputAccount.Email, inputAccount.Password)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		response, err := requests.AccessTokenRequest(account.ID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		err = json.NewDecoder(response.Body).Decode(&account)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("unable to parse response body"))
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(account)
	}
}
