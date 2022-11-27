package repositories

import (
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	Repo *Repository
}

func NewRoleRepository(repo *Repository) *RoleRepository {
	return &RoleRepository{Repo: repo}
}

func (repository *RoleRepository) Create(modelRole *models.Role) (err error) {
	if repository.Repo.Debug {
		err = repository.Repo.DB.Debug().Create(&modelRole).Error
	} else {
		err = repository.Repo.DB.Create(&modelRole).Error
	}

	return err
}

func (repository *RoleRepository) Update(modelRole *models.Role, data map[string]interface{}) (err error) {
	if repository.Repo.Debug {
		err = repository.Repo.DB.Debug().Model(&modelRole).Where("id = ?", modelRole.ID).Updates(data).Error
	} else {
		err = repository.Repo.DB.Model(&modelRole).Where("id = ?", modelRole.ID).Updates(data).Error
	}

	if err != nil {
		return err
	}

	if repository.Repo.Debug {
		err = repository.Repo.DB.Debug().Where("id = ?", modelRole.ID).Take(&modelRole).Error
	} else {
		err = repository.Repo.DB.Where("id = ?", modelRole.ID).Take(&modelRole).Error
	}

	return err
}

func (repository *RoleRepository) Delete(modelRole *models.Role) (int64, error) {
	var db *gorm.DB

	if repository.Repo.Debug {
		db = repository.Repo.DB.Debug().Where("id = ?", modelRole.ID).Delete(&models.Role{})
	} else {
		db = repository.Repo.DB.Where("id = ?", modelRole.ID).Delete(&models.Role{})
	}

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (repository *RoleRepository) FindAll() (modelsRole *[]models.Role, err error) {
	if repository.Repo.Debug {
		err = repository.Repo.DB.Debug().Find(&modelsRole).Error
	} else {
		err = repository.Repo.DB.Find(&modelsRole).Error
	}

	return modelsRole, err
}

func (repository *RoleRepository) FindOneByID(id uint) (modelRole *models.Role, err error) {
	if repository.Repo.Debug {
		err = repository.Repo.DB.Debug().First(&modelRole, id).Error
	} else {
		err = repository.Repo.DB.First(&modelRole, id).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return modelRole, fmt.Errorf("unable to find role %d: %s", id, err)
	}

	return modelRole, err
}

func (repository *RoleRepository) FindOneByName(roleName string) (modelRole *models.Role, err error) {
	if repository.Repo.Debug {
		err = repository.Repo.DB.Debug().Where("name = ?", roleName).First(&modelRole).Error
	} else {
		err = repository.Repo.DB.Where("name = ?", roleName).First(&modelRole).Error
	}

	return modelRole, err
}
