package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-account/models"
	"gorm.io/gorm"
	"net/http"
)

type RegisterHandler struct {
	db *gorm.DB
}

func NewRegisterHandler(db *gorm.DB) *RegisterHandler {
	return &RegisterHandler{db}
}

// type Handler interface
func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	roleRepository := models.RoleRepository{DB: h.db, Debug: false}
	var roles []models.Role
	role, err := roleRepository.FindOneByNameOrCreate("user")
	roles = append(roles, *role)

	accountRepository := models.AccountRepository{DB: h.db, Debug: false}
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
