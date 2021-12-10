package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-account/models"
	"gorm.io/gorm"
	"net/http"
)

func RegisterHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
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

		roleRepository := models.RoleRepository{DB: db, Debug: false}
		var roles []models.Role
		role, err := roleRepository.FindOneByNameOrCreate("user")
		roles = append(roles, *role)

		accountRepository := models.AccountRepository{DB: db, Debug: false}
		account, err := accountRepository.Create(&models.Account{
			Email:    inputAccount.Email,
			Password: inputAccount.Password,
			Roles:    roles,
		})

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(account)
	}
}
