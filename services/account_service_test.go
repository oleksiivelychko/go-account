package services

import (
	"github.com/oleksiivelychko/go-account/initdb"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/repositories"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	db, err := initdb.TestPrepare()
	if err != nil {
		t.Errorf("initialization test environment error: %s", err)
	}

	accountRepository := repositories.NewAccountRepository(db, false)
	accountService := NewAccountService(accountRepository)
	roleRepository := repositories.NewRoleRepository(db, false)
	roleService := NewRoleService(roleRepository)

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
	db, err := initdb.TestPrepare()
	if err != nil {
		t.Errorf("initialization test environment error: %s", err)
	}

	accountRepository := repositories.NewAccountRepository(db, false)
	accountService := NewAccountService(accountRepository)

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
	db, err := initdb.TestPrepare()
	if err != nil {
		t.Errorf("initialization test environment error: %s", err)
	}

	accountRepository := repositories.NewAccountRepository(db, false)
	accountService := NewAccountService(accountRepository)

	model, err := accountService.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	if err != nil {
		t.Errorf("unable to create account model: %s", err)
	}

	roleRepository := repositories.NewRoleRepository(db, false)
	roleService := NewRoleService(roleRepository)

	// assign roles to existing account
	var roles []models.Role

	roleManager, err := roleRepository.FindOneByName("manager")
	if err != nil {
		roleManager, err = roleService.Create(&models.Role{Name: "manager"})
	}
	roles = append(roles, *roleManager)

	roleSupplier, err := roleRepository.FindOneByName("supplier")
	if err != nil {
		roleSupplier, err = roleService.Create(&models.Role{Name: "supplier"})
	}
	roles = append(roles, *roleSupplier)

	model, err = accountRepository.AddRoles(model, roles)
	if err != nil {
		t.Errorf("unable to add roles to account model: %s", err)
	}

	if len(model.Roles) != 2 {
		t.Errorf("account model roles count mismatch: %d != 2", len(model.Roles))
	}
}

func TestDeleteRolesToAccount(t *testing.T) {
	db, err := initdb.TestPrepare()
	if err != nil {
		t.Errorf("initialization test environment error: %s", err)
	}

	accountRepository := repositories.NewAccountRepository(db, false)
	accountService := NewAccountService(accountRepository)

	model, err := accountService.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	if err != nil {
		t.Errorf("unable to create account model: %s", err)
	}

	roleRepository := repositories.NewRoleRepository(db, false)
	roleService := NewRoleService(roleRepository)

	// assign roles to existing account
	var roles []models.Role

	roleManager, err := roleService.Create(&models.Role{Name: "manager"})
	roles = append(roles, *roleManager)

	roleSupplier, err := roleService.Create(&models.Role{Name: "supplier"})
	roles = append(roles, *roleSupplier)

	model, err = accountRepository.DeleteRoles(model, roles)
	if err != nil {
		t.Errorf("unable to delete roles from account model: %s", err)
	}

	if len(model.Roles) != 0 {
		t.Errorf("account model roles count mismatch: %d != 0", len(model.Roles))
	}
}

func TestDeleteAccount(t *testing.T) {
	db, err := initdb.TestPrepare()
	if err != nil {
		t.Errorf("initialization test environment error: %s", err)
	}

	accountRepository := repositories.NewAccountRepository(db, false)
	accountService := NewAccountService(accountRepository)

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
