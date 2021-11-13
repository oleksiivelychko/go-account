package models

import (
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"testing"
)

func TestCreateRole(t *testing.T) {
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

	statement := "TRUNCATE roles RESTART IDENTITY CASCADE"
	sqlExec := orm.DB.Exec(statement)
	if sqlExec.Error != nil {
		t.Errorf("[sql exec `"+statement+"`] -> %s", sqlExec.Error)
	}

	roleRepository := RoleRepository{*orm}
	createdModel, err := roleRepository.Create(&Role{Name: "guest"})

	if err != nil {
		t.Errorf("[func (rr *RoleRepository) Create(model *Role) (*Role, error)] -> %s", err)
	}

	if createdModel.Name != "guest" {
		t.Errorf("[`Role` model.Name] -> %s != 'guest'", createdModel.Name)
	}
}

func TestUpdateRole(t *testing.T) {
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

	roleRepository := RoleRepository{*orm}
	model, err := roleRepository.FindOneByID(1)
	if err != nil {
		t.Errorf("[func (rr *RoleRepository) FindOneByID(uid uint) (*Role, error)] -> %s", err)
	}

	model.Name = "user"
	updatedModel, err := roleRepository.Update(model)
	if err != nil {
		t.Errorf("[func (rr *RoleRepository) Update(model *Role) (*Role, error)] -> %s", err)
	}

	if updatedModel.Name != "user" {
		t.Errorf("[`Role` model.Name] -> %s != 'user'", updatedModel.Name)
	}
}

func TestDelete(t *testing.T) {
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

	roleRepository := RoleRepository{*orm}
	model, err := roleRepository.FindOneByID(1)
	if err != nil {
		t.Errorf("[func (rr *RoleRepository) FindOneByID(uid uint) (*Account, error)] -> %s", err)
	}
	rowsAffected, err := roleRepository.Delete(model)
	if err != nil && rowsAffected == 0 {
		t.Errorf("[func (rr *RoleRepository) Delete(model *Role) (int64, error)] -> %s", err)
	}
}
