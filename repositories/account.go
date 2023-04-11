package repositories

import (
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"gorm.io/gorm"
)

type Account struct {
	Repo *Repository
}

func NewAccount(repo *Repository) *Account {
	return &Account{Repo: repo}
}

func (repo *Account) Create(modelAccount *models.Account) (err error) {
	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Create(&modelAccount).Error
	} else {
		err = repo.Repo.DB.Create(&modelAccount).Error
	}

	return err
}

func (repo *Account) Update(modelAccount *models.Account, data map[string]interface{}) (err error) {
	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Model(&modelAccount).Where("id = ?", modelAccount.ID).Updates(data).Error
	} else {
		err = repo.Repo.DB.Model(&modelAccount).Where("id = ?", modelAccount.ID).Updates(data).Error
	}

	if err != nil {
		return err
	}

	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Where("id = ?", modelAccount.ID).Take(&modelAccount).Error
	} else {
		err = repo.Repo.DB.Where("id = ?", modelAccount.ID).Take(&modelAccount).Error
	}

	return err
}

func (repo *Account) Delete(modelAccount *models.Account) (int64, error) {
	var err error = nil
	var db *gorm.DB

	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Model(&modelAccount).Association("Roles").Clear()
	} else {
		err = repo.Repo.DB.Model(&modelAccount).Association("Roles").Clear()
	}

	if err != nil {
		return 0, err
	}

	if repo.Repo.Debug {
		db = repo.Repo.DB.Debug().Where("id = ?", modelAccount.ID).Delete(&models.Account{})
	} else {
		db = repo.Repo.DB.Where("id = ?", modelAccount.ID).Delete(&models.Account{})
	}

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (repo *Account) FindAll() (modelsAccount *[]models.Account, err error) {
	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Preload("Roles").Find(&modelsAccount).Error
	} else {
		err = repo.Repo.DB.Preload("Roles").Find(&modelsAccount).Preload("Roles").Error
	}

	return modelsAccount, err
}

func (repo *Account) FindOneByID(id uint) (modelAccount *models.Account, err error) {
	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Preload("Roles").First(&modelAccount, id).Error
	} else {
		err = repo.Repo.DB.Preload("Roles").First(&modelAccount, id).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("unable to find account by ID %d: %s", id, err)
	}

	return modelAccount, err
}

func (repo *Account) FindOneByEmail(email string, withRoles bool) (*models.Account, error) {
	var db *gorm.DB
	var modelAccount *models.Account

	if repo.Repo.Debug {
		db = repo.Repo.DB.Debug().Where("email = ?", email)
	} else {
		db = repo.Repo.DB.Where("email = ?", email)
	}

	if withRoles {
		db = db.Preload("Roles").First(&modelAccount)
	} else {
		db = db.First(&modelAccount)
	}

	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("account not found: %s", db.Error)
	}

	return modelAccount, db.Error
}

func (repo *Account) AddRoles(modelAccount *models.Account, roles []*models.Role) (*models.Account, error) {
	var err error

	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Model(&modelAccount).Association("Roles").Append(&roles)
	} else {
		err = repo.Repo.DB.Model(&modelAccount).Association("Roles").Append(&roles)
	}

	if err != nil {
		return nil, err
	}

	return repo.FindOneByID(modelAccount.ID)
}

func (repo *Account) DeleteRoles(modelAccount *models.Account, roles []models.Role) (*models.Account, error) {
	var err error

	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Model(&modelAccount).Association("Roles").Delete(&roles)
	} else {
		err = repo.Repo.DB.Model(&modelAccount).Association("Roles").Delete(&roles)
	}

	if err != nil {
		return nil, err
	}

	return repo.FindOneByID(modelAccount.ID)
}
