package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/services"
	"net/http"
)

type Register struct {
	accountService *services.Account
	roleService    *services.Role
}

func NewRegister(accountService *services.Account, roleService *services.Role) *Register {
	return &Register{accountService, roleService}
}

func (handler *Register) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	resp.Header().Set("Content-Type", "application/json")

	var inputAccount models.Account

	err := json.NewDecoder(req.Body).Decode(&inputAccount)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		_, _ = resp.Write([]byte(err.Error()))
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
		resp.WriteHeader(http.StatusBadRequest)
		_, _ = resp.Write([]byte(err.Error()))
		return
	}

	resp.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(resp).Encode(handler.accountService.Serialize(account))
}
