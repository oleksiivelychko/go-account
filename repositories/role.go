package repositories

import (
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"gorm.io/gorm"
)

type Role struct {
	Repo *Repository
}

func NewRole(repo *Repository) *Role {
	return &Role{Repo: repo}
}

func (repo *Role) Create(modelRole *models.Role) (err error) {
	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Create(&modelRole).Error
	} else {
		err = repo.Repo.DB.Create(&modelRole).Error
	}

	return err
}

func (repo *Role) Update(modelRole *models.Role, data map[string]interface{}) (err error) {
	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Model(&modelRole).Where("id = ?", modelRole.ID).Updates(data).Error
	} else {
		err = repo.Repo.DB.Model(&modelRole).Where("id = ?", modelRole.ID).Updates(data).Error
	}

	if err != nil {
		return err
	}

	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Where("id = ?", modelRole.ID).Take(&modelRole).Error
	} else {
		err = repo.Repo.DB.Where("id = ?", modelRole.ID).Take(&modelRole).Error
	}

	return err
}

func (repo *Role) Delete(modelRole *models.Role) (int64, error) {
	var db *gorm.DB

	if repo.Repo.Debug {
		db = repo.Repo.DB.Debug().Where("id = ?", modelRole.ID).Delete(&models.Role{})
	} else {
		db = repo.Repo.DB.Where("id = ?", modelRole.ID).Delete(&models.Role{})
	}

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (repo *Role) FindAll() (modelsRole *[]models.Role, err error) {
	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Find(&modelsRole).Error
	} else {
		err = repo.Repo.DB.Find(&modelsRole).Error
	}

	return modelsRole, err
}

func (repo *Role) FindOneByID(id uint) (modelRole *models.Role, err error) {
	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().First(&modelRole, id).Error
	} else {
		err = repo.Repo.DB.First(&modelRole, id).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return modelRole, fmt.Errorf("unable to find role by ID %d: %s", id, err)
	}

	return modelRole, err
}

func (repo *Role) FindOneByName(roleName string) (modelRole *models.Role, err error) {
	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Where("name = ?", roleName).First(&modelRole).Error
	} else {
		err = repo.Repo.DB.Where("name = ?", roleName).First(&modelRole).Error
	}

	return modelRole, err
}
