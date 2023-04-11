package services

import (
	"errors"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/repositories"
	"gorm.io/gorm"
	"time"
)

type Role struct {
	repo *repositories.Role
}

func NewRole(rr *repositories.Role) *Role {
	return &Role{rr}
}

func (service *Role) Validate(role *models.Role) error {
	if role.Name != "" {
		roleFound, err := service.repo.FindOneByName(role.Name)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if roleFound.Name != "" {
			return errors.New("role already exists")
		}
	}

	return nil
}

func (service *Role) Create(role *models.Role) (*models.Role, error) {
	err := service.Validate(role)
	if err != nil {
		return &models.Role{}, err
	}

	err = service.repo.Create(role)
	if err != nil {
		return &models.Role{}, err
	}

	return role, nil
}

func (service *Role) Update(role *models.Role) (*models.Role, error) {
	err := service.Validate(role)
	if err != nil {
		return &models.Role{}, err
	}

	data := map[string]interface{}{
		"name":       role.Name,
		"updated_at": time.Now(),
	}

	err = service.repo.Update(role, data)
	if err != nil {
		return &models.Role{}, err
	}

	return role, nil
}

func (service *Role) Delete(role *models.Role) (int64, error) {
	return service.repo.Delete(role)
}

func (service *Role) FindOneByNameOrCreate(name string) (*models.Role, error) {
	role, err := service.repo.FindOneByName(name)
	if err != nil {
		role = &models.Role{Name: name}
		err = service.repo.Create(role)
	}

	return role, nil
}

func (service *Role) GetRepository() *repositories.Role {
	return service.repo
}
