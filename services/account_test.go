package services

import (
	"github.com/oleksiivelychko/go-account/db"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/repositories"
	"log"
	"testing"
)

func makeServices() (accountService *Account, roleService *Role) {
	sessionDB, err := db.PrepareTestDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	repo := repositories.NewRepository(sessionDB, false)
	accountRepo := repositories.NewAccount(repo)
	roleRepo := repositories.NewRole(repo)

	return NewAccount(accountRepo), NewRole(roleRepo)
}

func TestServices_CreateAccount(t *testing.T) {
	accountService, roleService := makeServices()

	var roles []models.Role
	role, _ := roleService.FindOneByNameOrCreate("user")
	roles = append(roles, *role)

	accountCreated, err := accountService.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
		Roles:    roles,
	})

	if err != nil {
		t.Fatal(err)
	}

	if accountCreated.Email != "test@test.test" {
		t.Errorf("email mismatch: %s != test@test.test", accountCreated.Email)
	}

	verifiedPasswordErr := accountService.VerifyPassword(accountCreated, "secret")
	if verifiedPasswordErr != nil {
		t.Error(verifiedPasswordErr)
	}

	if len(accountCreated.Roles) == 0 {
		t.Error("roles are missing")
	}

	if accountCreated.Roles[0].Name != "user" {
		t.Errorf("role name mismatch: %s != user", accountCreated.Roles[0].Name)
	}
}

func TestServices_UpdateAccount(t *testing.T) {
	accountService, _ := makeServices()

	account, err := accountService.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	if err != nil {
		t.Fatal(err)
	}

	account.Email = "test1@test1.test1"
	account.Password = "secret1"

	accountUpdated, err := accountService.Update(account)
	if err != nil {
		t.Fatal(err)
	}

	if accountUpdated.Email != "test1@test1.test1" {
		t.Errorf("email mismatch: %s != test1@test1.test1", accountUpdated.Email)
	}

	verifiedPasswordErr := accountService.VerifyPassword(accountUpdated, "secret1")
	if verifiedPasswordErr != nil {
		t.Error(verifiedPasswordErr)
	}
}

func TestServices_AddRolesToAccount(t *testing.T) {
	accountService, roleService := makeServices()

	account, err := accountService.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	if err != nil {
		t.Error(err)
	}

	// assign roles to existing account
	var roles []*models.Role

	roleManager, _ := roleService.GetRepository().FindOneByName("manager")
	roleManager, err = roleService.Create(&models.Role{Name: "manager"})
	roles = append(roles, roleManager)

	roleSupplier, _ := roleService.GetRepository().FindOneByName("supplier")
	roleSupplier, err = roleService.Create(&models.Role{Name: "supplier"})
	roles = append(roles, roleSupplier)

	accountWithRoles, err := accountService.GetRepository().AddRoles(account, roles)
	if err != nil {
		t.Fatal(err)
	}

	if len(accountWithRoles.Roles) != 2 {
		t.Errorf("roles count mismatch: %d != 2", len(accountWithRoles.Roles))
	}
}

func TestServices_DeleteRolesFromAccount(t *testing.T) {
	accountService, roleService := makeServices()

	account, err := accountService.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	if err != nil {
		t.Fatal(err)
	}

	// assign roles to existing account
	var roles []models.Role

	roleManager, err := roleService.Create(&models.Role{Name: "manager"})
	if err != nil {
		t.Fatal(err)
	}

	roles = append(roles, *roleManager)

	roleSupplier, err := roleService.Create(&models.Role{Name: "supplier"})
	if err != nil {
		t.Fatal(err)
	}

	roles = append(roles, *roleSupplier)

	account, err = accountService.GetRepository().DeleteRoles(account, roles)
	if err != nil {
		t.Fatal(err)
	}

	if len(account.Roles) != 0 {
		t.Errorf("roles count mismatch: %d != 0", len(account.Roles))
	}
}

func TestServices_DeleteAccount(t *testing.T) {
	accountService, _ := makeServices()

	account, err := accountService.Create(&models.Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	if err != nil {
		t.Fatal(err)
	}

	rowsAffected, err := accountService.Delete(account)
	if rowsAffected == 0 && err != nil {
		t.Error(err)
	}
}

func TestServices_FindAllAccountsWithRoles(t *testing.T) {
	accountService, roleService := makeServices()

	var roles []models.Role
	role, err := roleService.FindOneByNameOrCreate("user")
	if err != nil {
		t.Fatal(err)
	}

	roles = append(roles, *role)

	_, err = accountService.Create(&models.Account{
		Email:    "test1@test.test",
		Password: "secret",
		Roles:    roles,
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = accountService.Create(&models.Account{
		Email:    "test2@test.test",
		Password: "secret",
		Roles:    roles,
	})

	if err != nil {
		t.Fatal(err)
	}

	accounts, err := accountService.GetRepository().FindAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(*accounts) != 2 {
		t.Errorf("length mismatch: %d != 2", len(*accounts))
	}

	if cap(*accounts) != 20 {
		t.Errorf("capacity mismatch: %d != 20", cap(*accounts))
	}

	if len((*accounts)[0].Roles) != 1 {
		t.Errorf("roles length mismatch: %d != 1", len((*accounts)[0].Roles))
	}

	if cap((*accounts)[0].Roles) != 10 {
		t.Errorf("roles capacity mismatch: %d != 10", cap((*accounts)[0].Roles))
	}
}
