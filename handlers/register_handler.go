package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/services"
	"net/http"
)

type RegisterHandler struct {
	accountService *services.AccountService
	roleService    *services.RoleService
}

func NewRegisterHandler(a *services.AccountService, r *services.RoleService) *RegisterHandler {
	return &RegisterHandler{a, r}
}

func (handler *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var inputAccount models.Account
	err := json.NewDecoder(r.Body).Decode(&inputAccount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var roles []models.Role
	role, err := handler.roleService.FindOneByNameOrCreate("user")
	roles = append(roles, *role)

	account, err := handler.accountService.Create(&models.Account{
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
