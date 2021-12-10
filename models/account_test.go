package models

import (
	"github.com/oleksiivelychko/go-account/initdb"
	"gorm.io/gorm"
	"testing"
)

func initTest() (*gorm.DB, error) {
	initdb.LoadEnv()
	db, err := initdb.TestDB()
	err = AutoMigrate(db)

	statement := "TRUNCATE accounts RESTART IDENTITY CASCADE"
	sqlExec := db.Exec(statement)
	if sqlExec.Error != nil {
		return nil, sqlExec.Error
	}

	statement = "TRUNCATE roles RESTART IDENTITY CASCADE"
	sqlExec = db.Exec(statement)
	if sqlExec.Error != nil {
		return nil, sqlExec.Error
	}

	return db, err
}

func TestCreateAccount(t *testing.T) {
	db, err := initTest()
	if err != nil {
		t.Errorf("error during initialization test environment: %s", err)
	}

	roleRepository := RoleRepository{db, false}
	var roles []Role
	role, _ := roleRepository.FindOneByNameOrCreate("user")
	roles = append(roles, *role)

	accountRepository := AccountRepository{db, false}
	createdModel, err := accountRepository.Create(&Account{
		Email:    "test@test.test",
		Password: "secret",
		Roles:    roles,
	})

	if err != nil {
		t.Errorf("[func (ar *AccountRepository) Create(model *Account) (*Account, error)] -> %s", err)
	}

	if createdModel.Email != "test@test.test" {
		t.Errorf("[`Account` model.Email] -> %s != 'test@test.test'", createdModel.Email)
	}

	verifiedPassword := createdModel.VerifyPassword("secret")
	if verifiedPassword != nil {
		t.Errorf("[func (model *Account) VerifyPassword(password string) error] -> %s", verifiedPassword)
	}

	if len(createdModel.Roles) == 0 {
		t.Errorf("[`Account` model.Roles len] -> %d == '0'", 0)
	}

	if createdModel.Roles[0].Name != "user" {
		t.Errorf("[`Account` model.Roles 'user'] -> %s != 'user'", createdModel.Roles[0].Name)
	}
}

func TestUpdateAccount(t *testing.T) {
	db, err := initTest()
	if err != nil {
		t.Errorf("error during initialization test environment: %s", err)
	}

	accountRepository := AccountRepository{db, false}
	model, err := accountRepository.Create(&Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	if err != nil {
		t.Errorf("[func (ar *AccountRepository) Create(model *Account) (*Account, error)] -> %s", err)
	}

	model.Email = "test1@test1.test1"
	model.Password = "secret1"
	updatedModel, err := accountRepository.Update(model)
	if err != nil {
		t.Errorf("[func (ar *AccountRepository) Update(model *Account) (*Account, error)] -> %s", err)
	}

	if updatedModel.Email != "test1@test1.test1" {
		t.Errorf("[`Account` model.Email] -> %s != 'test1@test1.test1'", updatedModel.Email)
	}

	verifiedPassword := updatedModel.VerifyPassword("secret1")
	if verifiedPassword != nil {
		t.Errorf("[func (model *Account) VerifyPassword(password string) error] -> %s", verifiedPassword)
	}
}

func TestAddRolesToAccount(t *testing.T) {
	db, err := initTest()
	if err != nil {
		t.Errorf("error during initialization test environment: %s", err)
	}

	accountRepository := AccountRepository{db, false}
	model, err := accountRepository.Create(&Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	if err != nil {
		t.Errorf("[func (ar *AccountRepository) Create(model *Account) (*Account, error)] -> %s", err)
	}

	// assign roles to exists account
	roleRepository := RoleRepository{db, false}
	var roles []Role
	roleManager, err := roleRepository.FindOneByName("manager")
	if err != nil {
		roleManager, err = roleRepository.Create(&Role{Name: "manager"})
	}
	roles = append(roles, *roleManager)

	roleSupplier, err := roleRepository.FindOneByName("supplier")
	if err != nil {
		roleSupplier, err = roleRepository.Create(&Role{Name: "supplier"})
	}
	roles = append(roles, *roleSupplier)

	model, err = accountRepository.AddRoles(model, roles)
	if err != nil {
		t.Errorf("[`Account` model.AddRoles] -> %s", err)
	}
	if len(model.Roles) != 2 {
		t.Errorf("[`Account` model.Roles len] -> %d != 2", len(model.Roles))
	}
}

func TestDeleteRolesToAccount(t *testing.T) {
	db, err := initTest()
	if err != nil {
		t.Errorf("error during initialization test environment: %s", err)
	}

	accountRepository := AccountRepository{db, false}
	model, err := accountRepository.Create(&Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	if err != nil {
		t.Errorf("[func (ar *AccountRepository) Create(model *Account) (*Account, error)] -> %s", err)
	}

	// assign roles to exists account
	roleRepository := RoleRepository{db, false}
	var roles []Role
	roleManager, err := roleRepository.Create(&Role{Name: "manager"})
	roles = append(roles, *roleManager)
	roleSupplier, err := roleRepository.Create(&Role{Name: "supplier"})
	roles = append(roles, *roleSupplier)

	model, err = accountRepository.DeleteRoles(model, roles)
	if err != nil {
		t.Errorf("[`Account` model.AddRoles] -> %s", err)
	}
	if len(model.Roles) != 0 {
		t.Errorf("[`Account` model.Roles len] -> %d != 0", len(model.Roles))
	}
}

func TestDeleteAccount(t *testing.T) {
	db, err := initTest()
	if err != nil {
		t.Errorf("error during initialization test environment: %s", err)
	}

	accountRepository := AccountRepository{db, false}
	model, err := accountRepository.Create(&Account{
		Email:    "test@test.test",
		Password: "secret",
	})

	if err != nil {
		t.Errorf("[func (ar *AccountRepository) Create(model *Account) (*Account, error)] -> %s", err)
	}

	rowsAffected, err := accountRepository.Delete(model)
	if rowsAffected == 0 && err != nil {
		t.Errorf("[func (ar *AccountRepository) Delete(model *Account) error] -> %s", err)
	}
}
