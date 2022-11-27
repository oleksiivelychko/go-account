package services

import (
	"github.com/oleksiivelychko/go-account/models"
	"testing"
)

func TestCreateRole(t *testing.T) {
	_, roleService := InitTest()

	createdRole, err := roleService.Create(&models.Role{Name: "guest"})

	if err != nil {
		t.Errorf("unable to create role model: %s", err)
	}

	if createdRole.Name != "guest" {
		t.Errorf("role model name mismatch: '%s' != 'guest'", createdRole.Name)
	}
}

func TestUpdateRole(t *testing.T) {
	_, roleService := InitTest()

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
	_, roleService := InitTest()

	role, err := roleService.Create(&models.Role{Name: "guest"})

	rowsAffected, err := roleService.Delete(role)
	if err != nil && rowsAffected == 0 {
		t.Errorf("unable to delete role model: %s", err)
	}
}

func TestFindAllRoles(t *testing.T) {
	_, roleService := InitTest()

	var roles []models.Role
	role, _ := roleService.FindOneByNameOrCreate("user")
	roles = append(roles, *role)

	_, err := roleService.Create(&models.Role{Name: "guest"})
	_, err = roleService.Create(&models.Role{Name: "user"})

	modelsRole, err := roleService.GetRepository().FindAll()
	if err != nil {
		t.Errorf("unable to find all role models: %s", err)
	}

	if len(*modelsRole) != 2 {
		t.Errorf("role models length mismatch: %d != 2", len(*modelsRole))
	}

	if cap(*modelsRole) != 20 {
		t.Errorf("role models capacity mismatch: %d != 20", cap(*modelsRole))
	}
}
