package models

import (
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	err := godotenv.Load("./../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	orm := InitDB(true)
	DB, err := orm.DB.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Println(err)
		}
	}(DB)

	statement := "TRUNCATE accounts RESTART IDENTITY CASCADE"
	sqlExec := orm.DB.Exec(statement)
	if sqlExec.Error != nil {
		t.Errorf("[sql exec `"+statement+"`] -> %s", sqlExec.Error)
	}

	roleRepository := RoleRepository{*orm}
	var roles []Role
	roleAdmin, err := roleRepository.FindOneByName("admin")
	if err != nil {
		roleAdmin, err = roleRepository.Create(&Role{Name: "admin"})
	}
	roles = append(roles, *roleAdmin)

	accountRepository := AccountRepository{*orm}
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

	if createdModel.Roles[0].Name != "admin" {
		t.Errorf("[`Account` model.Roles 'admin'] -> %s != 'admin'", createdModel.Roles[0].Name)
	}
}

func TestUpdateAccount(t *testing.T) {
	err := godotenv.Load("./../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	orm := InitDB(true)
	DB, err := orm.DB.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Println(err)
		}
	}(DB)

	accountRepository := AccountRepository{*orm}
	model, err := accountRepository.FindOneByID(1)
	if err != nil {
		t.Errorf("[func (ar *AccountRepository) FindOneByID(uid uint) (*Account, error)] -> %s", err)
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
	err := godotenv.Load("./../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	orm := InitDB(true)
	DB, err := orm.DB.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Println(err)
		}
	}(DB)

	accountRepository := AccountRepository{*orm}
	model, err := accountRepository.FindOneByID(1)
	if err != nil {
		t.Errorf("[func (ar *AccountRepository) FindOneByID(uid uint) (*Account, error)] -> %s", err)
	}

	// Assign roles to exists account
	roleRepository := RoleRepository{*orm}
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
	if len(model.Roles) != 3 {
		t.Errorf("[`Account` model.Roles len] -> %d != '3'", len(model.Roles))
	}
}

func TestDeleteRolesToAccount(t *testing.T) {
	err := godotenv.Load("./../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	orm := InitDB(true)
	DB, err := orm.DB.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Println(err)
		}
	}(DB)

	accountRepository := AccountRepository{*orm}
	model, err := accountRepository.FindOneByID(1)
	if err != nil {
		t.Errorf("[func (ar *AccountRepository) FindOneByID(uid uint) (*Account, error)] -> %s", err)
	}

	// Assign roles to exists account
	roleRepository := RoleRepository{*orm}
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

	model, err = accountRepository.DeleteRoles(model, roles)
	if err != nil {
		t.Errorf("[`Account` model.AddRoles] -> %s", err)
	}
	if len(model.Roles) != 1 {
		t.Errorf("[`Account` model.Roles len] -> %d != '1'", len(model.Roles))
	}
}

func TestDeleteAccount(t *testing.T) {
	err := godotenv.Load("./../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	orm := InitDB(true)
	DB, err := orm.DB.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Println(err)
		}
	}(DB)

	accountRepository := AccountRepository{*orm}
	model, err := accountRepository.FindOneByID(1)
	if err != nil {
		t.Errorf("[func (ar *AccountRepository) FindOneByID(uid uint) (*Account, error)] -> %s", err)
	}

	rowsAffected, err := accountRepository.Delete(model)
	if rowsAffected == 0 && err != nil {
		t.Errorf("[func (ar *AccountRepository) Delete(model *Account) error] -> %s", err)
	}
}
