package models

import (
	"testing"
)

func TestCreateRole(t *testing.T) {
	db, err := initTest()
	if err != nil {
		t.Errorf("error during initialization test environment: %s", err)
	}

	roleRepository := RoleRepository{db, false}
	createdRole, err := roleRepository.Create(&Role{Name: "guest"})

	if err != nil {
		t.Errorf("[func (rr *RoleRepository) Create(model *Role) (*Role, error)] -> %s", err)
	}

	if createdRole.Name != "guest" {
		t.Errorf("[`Role` model.Name] -> %s != 'guest'", createdRole.Name)
	}
}

func TestUpdateRole(t *testing.T) {
	db, err := initTest()
	if err != nil {
		t.Errorf("error during initialization test environment: %s", err)
	}

	roleRepository := RoleRepository{db, false}
	existsRole, err := roleRepository.Create(&Role{Name: "guest"})

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
	db, err := initTest()
	if err != nil {
		t.Errorf("error during initialization test environment: %s", err)
	}

	roleRepository := RoleRepository{db, false}
	role, err := roleRepository.Create(&Role{Name: "guest"})

	rowsAffected, err := roleRepository.Delete(role)
	if err != nil && rowsAffected == 0 {
		t.Errorf("[func (rr *RoleRepository) Delete(model *Role) (int64, error)] -> %s", err)
	}
}
