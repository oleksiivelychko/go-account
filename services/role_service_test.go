package services

import (
	"github.com/oleksiivelychko/go-account/initdb"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/repositories"
	"testing"
)

func TestCreateRole(t *testing.T) {
	db, err := initdb.TestPrepare()
	if err != nil {
		t.Errorf("initialization test environment error: %s", err)
	}

	roleRepository := repositories.NewRoleRepository(db, false)
	roleService := NewRoleService(roleRepository)

	createdRole, err := roleService.Create(&models.Role{Name: "guest"})

	if err != nil {
		t.Errorf("unable to create role model: %s", err)
	}

	if createdRole.Name != "guest" {
		t.Errorf("role model name mismatch: '%s' != 'guest'", createdRole.Name)
	}
}

func TestUpdateRole(t *testing.T) {
	db, err := initdb.TestPrepare()
	if err != nil {
		t.Errorf("initialization test environment error: %s", err)
	}

	roleRepository := repositories.NewRoleRepository(db, false)
	roleService := NewRoleService(roleRepository)

	role, err := roleService.Create(&models.Role{Name: "guest"})
	role.Name = "user"

	updatedRole, err := roleService.Update(role)
	if err != nil {
		t.Errorf("unable to update role model: %s", err)
	}

	if updatedRole.Name != "user" {
		t.Errorf("role name mismatch: '%s' != 'user'", updatedRole.Name)
	}
}

func TestDelete(t *testing.T) {
	db, err := initdb.TestPrepare()
	if err != nil {
		t.Errorf("initialization test environment error: %s", err)
	}

	roleRepository := repositories.NewRoleRepository(db, false)
	roleService := NewRoleService(roleRepository)

	role, err := roleService.Create(&models.Role{Name: "guest"})

	rowsAffected, err := roleRepository.Delete(role)
	if err != nil && rowsAffected == 0 {
		t.Errorf("unable to delete role model: %s", err)
	}
}
