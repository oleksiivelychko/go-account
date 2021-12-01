package models

import (
	"database/sql"
	"github.com/oleksiivelychko/go-account/initdb"
	"log"
	"testing"
)

func TestCreateRole(t *testing.T) {
	initdb.LoadEnv()
	db, err := initdb.TestDB()
	dbConnection, _ := db.DB()

	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Println(err)
		}
	}(dbConnection)

	statement := "TRUNCATE roles RESTART IDENTITY CASCADE"
	sqlExec := db.Exec(statement)
	if sqlExec.Error != nil {
		t.Errorf("[sql exec `"+statement+"`] -> %s", sqlExec.Error)
	}

	roleRepository := RoleRepository{db, true}
	createdRole, err := roleRepository.Create(&Role{Name: "guest"})

	if err != nil {
		t.Errorf("[func (rr *RoleRepository) Create(model *Role) (*Role, error)] -> %s", err)
	}

	if createdRole.Name != "guest" {
		t.Errorf("[`Role` model.Name] -> %s != 'guest'", createdRole.Name)
	}
}

func TestUpdateRole(t *testing.T) {
	initdb.LoadEnv()
	db, err := initdb.TestDB()
	dbConnection, _ := db.DB()

	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Println(err)
		}
	}(dbConnection)

	roleRepository := RoleRepository{db, true}
	existsRole, err := roleRepository.FindOneByID(1)
	if err != nil {
		t.Errorf("[func (rr *RoleRepository) FindOneByID(uid uint) (*Role, error)] -> %s", err)
	}

	existsRole.Name = "user"
	updatedRole, err := roleRepository.Update(existsRole)
	if err != nil {
		t.Errorf("[func (rr *RoleRepository) Update(model *Role) (*Role, error)] -> %s", err)
	}

	if updatedRole.Name != "user" {
		t.Errorf("[`Role` model.Name] -> %s != 'user'", updatedRole.Name)
	}
}

func TestDelete(t *testing.T) {
	initdb.LoadEnv()
	db, err := initdb.TestDB()
	dbConnection, _ := db.DB()

	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Println(err)
		}
	}(dbConnection)

	roleRepository := RoleRepository{db, true}
	role, err := roleRepository.FindOneByID(1)
	if err != nil {
		t.Errorf("[func (rr *RoleRepository) FindOneByID(uid uint) (*Account, error)] -> %s", err)
	}
	rowsAffected, err := roleRepository.Delete(role)
	if err != nil && rowsAffected == 0 {
		t.Errorf("[func (rr *RoleRepository) Delete(model *Role) (int64, error)] -> %s", err)
	}
}
