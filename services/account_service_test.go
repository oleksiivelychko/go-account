package services

import (
	"fmt"
	"github.com/oleksiivelychko/go-account/initdb"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/repositories"
	"os"
	"testing"
)

func InitTest() (accountService *AccountService, roleService *RoleService) {
	db, err := initdb.TestPrepare()
	if err != nil {
		fmt.Printf("initialization test environment error: %s", err)
		os.Exit(0)
	}

	repository := repositories.NewRepository(db, false)
	accountRepository := repositories.NewAccountRepository(repository)
	roleRepository := repositories.NewRoleRepository(repository)

	return NewAccountService(accountRepository), NewRoleService(roleRepository)
}

func TestCreateAccount(t *testing.T) {
	accountService, roleService := InitTest()

	var roles []models.Role
	role, _ := roleService.FindOneByNameOrCreate("user")
	roles = append(roles, *role)

	createdModel, err := accountService.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
		Roles:    roles,
	})

	if err != nil {
		t.Errorf("unable to create account model: %s", err)
	}

	if createdModel.Email != "test@test.test" {
		t.Errorf("account email mismatch: '%s' != 'test@test.test'", createdModel.Email)
	}

	verifiedPasswordErr := accountService.VerifyPassword(createdModel, "secret")
	if verifiedPasswordErr != nil {
		t.Errorf("unable to verify account password: %s", verifiedPasswordErr)
	}

	if len(createdModel.Roles) == 0 {
		t.Error("account model doesn't have any role")
	}

	if createdModel.Roles[0].Name != "user" {
		t.Errorf("account role mismatch: '%s' != 'user'", createdModel.Roles[0].Name)
	}
}

func TestUpdateAccount(t *testing.T) {
	accountService, _ := InitTest()

	model, err := accountService.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	if err != nil {
		t.Errorf("unable to create account model: %s", err)
	}

	model.Email = "test1@test1.test1"
	model.Password = "secret1"

	updatedModel, err := accountService.Update(model)
	if err != nil {
		t.Errorf("unable to update account model: %s", err)
	}

	if updatedModel.Email != "test1@test1.test1" {
		t.Errorf("account email mismatch: '%s' != 'test1@test1.test1'", updatedModel.Email)
	}

	verifiedPasswordErr := accountService.VerifyPassword(updatedModel, "secret1")
	if verifiedPasswordErr != nil {
		t.Errorf("unable to verify account password: %s", verifiedPasswordErr)
	}
}

func TestAddRolesToAccount(t *testing.T) {
	accountService, roleService := InitTest()

	modelAccount, err := accountService.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	if err != nil {
		t.Errorf("unable to create account model: %s", err)
	}

	// assign roles to existing account
	var roles []*models.Role

	roleManager, err := roleService.GetRepository().FindOneByName("manager")
	if err != nil {
		roleManager, err = roleService.Create(&models.Role{Name: "manager"})
	}
	roles = append(roles, roleManager)

	roleSupplier, err := roleService.GetRepository().FindOneByName("supplier")
	if err != nil {
		roleSupplier, err = roleService.Create(&models.Role{Name: "supplier"})
	}
	roles = append(roles, roleSupplier)

	modelAccountWithRoles, err := accountService.GetRepository().AddRoles(modelAccount, roles)
	if err != nil {
		t.Errorf("unable to add roles to account model: %s", err)
	}

	if len(modelAccountWithRoles.Roles) != 2 {
		t.Errorf("account model roles count mismatch: %d != 2", len(modelAccountWithRoles.Roles))
	}
}

func TestDeleteRolesToAccount(t *testing.T) {
	accountService, roleService := InitTest()

	model, err := accountService.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	if err != nil {
		t.Errorf("unable to create account model: %s", err)
	}

	// assign roles to existing account
	var roles []models.Role

	roleManager, err := roleService.Create(&models.Role{Name: "manager"})
	roles = append(roles, *roleManager)

	roleSupplier, err := roleService.Create(&models.Role{Name: "supplier"})
	roles = append(roles, *roleSupplier)

	model, err = accountService.GetRepository().DeleteRoles(model, roles)
	if err != nil {
		t.Errorf("unable to delete roles from account model: %s", err)
	}

	if len(model.Roles) != 0 {
		t.Errorf("account model roles count mismatch: %d != 0", len(model.Roles))
	}
}

func TestDeleteAccount(t *testing.T) {
	accountService, _ := InitTest()

	model, err := accountService.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	if err != nil {
		t.Errorf("unable to create account model: %s", err)
	}

	rowsAffected, err := accountService.Delete(model)
	if rowsAffected == 0 && err != nil {
		t.Errorf("unable to delete account model: %s", err)
	}
}

func TestFindAllAccountsWithRoles(t *testing.T) {
	accountService, roleService := InitTest()

	var roles []models.Role
	role, _ := roleService.FindOneByNameOrCreate("user")
	roles = append(roles, *role)

	_, err := accountService.Create(&models.Account{
		Email:    "test1@test.test",
		Password: "secret",
		Roles:    roles,
	})

	_, err = accountService.Create(&models.Account{
		Email:    "test2@test.test",
		Password: "secret",
		Roles:    roles,
	})

	modelsAccount, err := accountService.GetRepository().FindAll()
	if err != nil {
		t.Errorf("unable to find all account models: %s", err)
	}

	if len(*modelsAccount) != 2 {
		t.Errorf("account models length mismatch: %d != 2", len(*modelsAccount))
	}

	if cap(*modelsAccount) != 20 {
		t.Errorf("account models capacity mismatch: %d != 20", cap(*modelsAccount))
	}

	if len((*modelsAccount)[0].Roles) != 1 {
		t.Errorf("account model roles length mismatch: %d != 1", len((*modelsAccount)[0].Roles))
	}

	if cap((*modelsAccount)[0].Roles) != 10 {
		t.Errorf("account model roles capacity mismatch: %d != 10", cap((*modelsAccount)[0].Roles))
	}
}
