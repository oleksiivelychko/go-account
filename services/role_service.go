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

func NewRoleService(r *repositories.RoleRepository) *RoleService {
	return &RoleService{r}
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

func (service *RoleService) Create(role *models.Role) (*models.Role, error) {
	err := service.Validate(role)
	if err != nil {
		return &models.Role{}, err
	}

	return service.repository.Create(role)
}

func (service *RoleService) Update(role *models.Role) (*models.Role, error) {
	err := service.Validate(role)
	if err != nil {
		return &models.Role{}, err
	}

	data := map[string]interface{}{
		"name":       role.Name,
		"updated_at": time.Now(),
	}

	return service.repository.Update(role, data)
}

func (service *RoleService) Delete(role *models.Role) (int64, error) {
	return service.repository.Delete(role)
}

func (service *RoleService) FindOneByNameOrCreate(name string) (*models.Role, error) {
	role, err := service.repository.FindOneByName(name)
	if err != nil {
		role, err = service.repository.Create(&models.Role{Name: name})
	}

	return role, nil
}

func (service *RoleService) GetRepository() *repositories.RoleRepository {
	return service.repository
}
