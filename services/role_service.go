package services

import (
	"errors"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/repositories"
	"gorm.io/gorm"
	"time"
)

type RoleService struct {
	repository *repositories.RoleRepository
}

func NewRoleService(rr *repositories.RoleRepository) *RoleService {
	return &RoleService{rr}
}

func (service *RoleService) Validate(role *models.Role) error {
	if role.Name != "" {
		existsRole, err := service.repository.FindOneByName(role.Name)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if existsRole.Name != "" {
			return errors.New("role already exists")
		}
	}

	return nil
}

func (service *RoleService) Create(modelRole *models.Role) (*models.Role, error) {
	err := service.Validate(modelRole)
	if err != nil {
		return &models.Role{}, err
	}

	err = service.repository.Create(modelRole)
	if err != nil {
		return &models.Role{}, err
	}

	return modelRole, nil
}

func (service *RoleService) Update(modelRole *models.Role) (*models.Role, error) {
	err := service.Validate(modelRole)
	if err != nil {
		return &models.Role{}, err
	}

	data := map[string]interface{}{
		"name":       modelRole.Name,
		"updated_at": time.Now(),
	}

	err = service.repository.Update(modelRole, data)
	if err != nil {
		return &models.Role{}, err
	}

	return modelRole, nil
}

func (service *RoleService) Delete(modelRole *models.Role) (int64, error) {
	return service.repository.Delete(modelRole)
}

func (service *RoleService) FindOneByNameOrCreate(roleName string) (*models.Role, error) {
	modelRole, err := service.repository.FindOneByName(roleName)
	if err != nil {
		modelRole = &models.Role{Name: roleName}
		err = service.repository.Create(modelRole)
	}

	return modelRole, nil
}

func (service *RoleService) GetRepository() *repositories.RoleRepository {
	return service.repository
}
