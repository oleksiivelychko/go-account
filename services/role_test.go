package services

import (
	"github.com/oleksiivelychko/go-account/db"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/repositories"
	"log"
	"testing"
)

func makeService() *Role {
	sessionDB, err := db.PrepareTestDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	repo := repositories.NewRepository(sessionDB, false)
	roleRepo := repositories.NewRole(repo)

	return NewRole(roleRepo)
}

func TestServices_CreateRole(t *testing.T) {
	service := makeService()
	role, err := service.Create(&models.Role{Name: "guest"})

	if err != nil {
		t.Fatal(err)
	}

	if role.Name != "guest" {
		t.Errorf("name mismatch: %s != guest", role.Name)
	}
}

func TestServices_UpdateRole(t *testing.T) {
	service := makeService()
	role, err := service.Create(&models.Role{Name: "guest"})
	role.Name = "user"

	roleUpdated, err := service.Update(role)
	if err != nil {
		t.Fatal(err)
	}

	if roleUpdated.Name != "user" {
		t.Errorf("name mismatch: %s != user", roleUpdated.Name)
	}
}

func TestServices_DeleteRole(t *testing.T) {
	service := makeService()
	role, err := service.Create(&models.Role{Name: "guest"})

	rowsAffected, err := service.Delete(role)
	if err != nil && rowsAffected == 0 {
		t.Error(err)
	}
}

func TestServices_FindOrCreateRole(t *testing.T) {
	service := makeService()

	_, err := service.FindOneByNameOrCreate("user")
	if err != nil {
		t.Error(err)
	}
}

func TestServices_FindAllRoles(t *testing.T) {
	service := makeService()

	_, err := service.Create(&models.Role{Name: "guest"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.Create(&models.Role{Name: "user"})
	if err != nil {
		t.Fatal(err)
	}

	roles, err := service.GetRepository().FindAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(*roles) != 2 {
		t.Errorf("length mismatch: %d != 2", len(*roles))
	}

	if cap(*roles) != 20 {
		t.Errorf("capacity mismatch: %d != 20", cap(*roles))
	}
}
