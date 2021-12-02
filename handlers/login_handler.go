package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-account/models"
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

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(account)
	}
}
