package repositories

import (
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	DB    *gorm.DB
	Debug bool
}

func NewRoleRepository(db *gorm.DB, debug bool) *RoleRepository {
	return &RoleRepository{db, debug}
}

func (repository *RoleRepository) Create(role *models.Role) (*models.Role, error) {
	var err error

	if repository.Debug {
		err = repository.DB.Debug().Create(&role).Error
	} else {
		err = repository.DB.Create(&role).Error
	}

	if err != nil {
		return &models.Role{}, err
	}

	return role, nil
}

func (repository *RoleRepository) Update(role *models.Role, data map[string]interface{}) (*models.Role, error) {
	var err error

	if repository.Debug {
		err = repository.DB.Debug().Model(&role).Where("id = ?", role.ID).Updates(data).Error
	} else {
		err = repository.DB.Model(&role).Where("id = ?", role.ID).Updates(data).Error
	}

	if err != nil {
		return role, err
	}

	if repository.Debug {
		err = repository.DB.Debug().Where("id = ?", role.ID).Take(&role).Error
	} else {
		err = repository.DB.Where("id = ?", role.ID).Take(&role).Error
	}

	if err != nil {
		return role, err
	}

	return role, nil
}

func (repository *RoleRepository) Delete(role *models.Role) (int64, error) {
	var db *gorm.DB
	if repository.Debug {
		db = repository.DB.Debug().Where("id = ?", role.ID).Delete(&models.Role{})
	} else {
		db = repository.DB.Where("id = ?", role.ID).Delete(&models.Role{})
	}

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (repository *RoleRepository) FindAll() (roles *[]models.Role, err error) {
	if repository.Debug {
		err = repository.DB.Debug().Find(&roles).Error
	} else {
		err = repository.DB.Find(&roles).Error
	}

	if err != nil {
		return &[]models.Role{}, err
	}
	return roles, err
}

func (repository *RoleRepository) FindOneByID(uid uint) (role *models.Role, err error) {
	var db *gorm.DB

	if repository.Debug {
		db = repository.DB.Debug().First(&role, uid)
	} else {
		db = repository.DB.First(&role, uid)
	}

	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		return role, fmt.Errorf("role `%d` doesn't managed to find: %s", uid, db.Error)
	}

	return role, db.Error
}

func (repository *RoleRepository) FindOneByName(roleName string) (*models.Role, error) {
	var db *gorm.DB
	var role = &models.Role{}

	if repository.Debug {
		db = repository.DB.Debug().Where("name = ?", roleName).First(role)
	} else {
		db = repository.DB.Where("name = ?", roleName).First(role)
	}

	return role, db.Error
}
