package repositories

import (
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"gorm.io/gorm"
)

type AccountRepository struct {
	Repo *Repository
}

func NewAccountRepository(repo *Repository) *AccountRepository {
	return &AccountRepository{Repo: repo}
}

func (ar *AccountRepository) Create(modelAccount *models.Account) (err error) {
	if ar.Repo.Debug {
		err = ar.Repo.DB.Debug().Create(&modelAccount).Error
	} else {
		err = ar.Repo.DB.Create(&modelAccount).Error
	}

	return err
}

func (ar *AccountRepository) Update(modelAccount *models.Account, data map[string]interface{}) (err error) {
	if ar.Repo.Debug {
		err = ar.Repo.DB.Debug().Model(&modelAccount).Where("id = ?", modelAccount.ID).Updates(data).Error
	} else {
		err = ar.Repo.DB.Model(&modelAccount).Where("id = ?", modelAccount.ID).Updates(data).Error
	}

	if err != nil {
		return err
	}

	if ar.Repo.Debug {
		err = ar.Repo.DB.Debug().Where("id = ?", modelAccount.ID).Take(&modelAccount).Error
	} else {
		err = ar.Repo.DB.Where("id = ?", modelAccount.ID).Take(&modelAccount).Error
	}

	return err
}

func (ar *AccountRepository) Delete(modelAccount *models.Account) (int64, error) {
	var err error = nil
	var db *gorm.DB

	if ar.Repo.Debug {
		err = ar.Repo.DB.Debug().Model(&modelAccount).Association("Roles").Clear()
	} else {
		err = ar.Repo.DB.Model(&modelAccount).Association("Roles").Clear()
	}

	if err != nil {
		return 0, err
	}

	if ar.Repo.Debug {
		db = ar.Repo.DB.Debug().Where("id = ?", modelAccount.ID).Delete(&models.Account{})
	} else {
		db = ar.Repo.DB.Where("id = ?", modelAccount.ID).Delete(&models.Account{})
	}

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (ar *AccountRepository) FindAll() (modelsAccount *[]models.Account, err error) {
	if ar.Repo.Debug {
		err = ar.Repo.DB.Debug().Find(&modelsAccount).Error
	} else {
		err = ar.Repo.DB.Find(&modelsAccount).Error
	}

	return modelsAccount, err
}

func (ar *AccountRepository) FindOneByID(id uint) (modelAccount *models.Account, err error) {
	if ar.Repo.Debug {
		err = ar.Repo.DB.Debug().Preload("Roles").First(&modelAccount, id).Error
	} else {
		err = ar.Repo.DB.Preload("Roles").First(&modelAccount, id).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("account not found: %s", err)
	}

	return modelAccount, err
}

func (ar *AccountRepository) FindOneByEmail(email string, withRoles bool) (*models.Account, error) {
	var db *gorm.DB
	var modelAccount *models.Account

	if ar.Repo.Debug {
		db = ar.Repo.DB.Debug().Where("email = ?", email)
	} else {
		db = ar.Repo.DB.Where("email = ?", email)
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

func (ar *AccountRepository) AddRoles(modelAccount *models.Account, roles []models.Role) (*models.Account, error) {
	var err error

	if ar.Repo.Debug {
		err = ar.Repo.DB.Debug().Model(&modelAccount).Association("Roles").Append(&roles)
	} else {
		err = ar.Repo.DB.Model(&modelAccount).Association("Roles").Append(&roles)
	}

	if err != nil {
		return nil, err
	}

	return ar.FindOneByID(modelAccount.ID)
}

func (ar *AccountRepository) DeleteRoles(modelAccount *models.Account, roles []models.Role) (*models.Account, error) {
	var err error

	if ar.Repo.Debug {
		err = ar.Repo.DB.Debug().Model(&modelAccount).Association("Roles").Delete(&roles)
	} else {
		err = ar.Repo.DB.Model(&modelAccount).Association("Roles").Delete(&roles)
	}

	if err != nil {
		return nil, err
	}

	return ar.FindOneByID(modelAccount.ID)
}
