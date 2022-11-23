package repositories

import (
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"gorm.io/gorm"
)

type AccountRepository struct {
	DB    *gorm.DB
	Debug bool
}

func NewAccountRepository(db *gorm.DB, debug bool) *AccountRepository {
	return &AccountRepository{db, debug}
}

func (repository *AccountRepository) Create(account *models.Account) (*models.Account, error) {
	var err error

	if repository.Debug {
		err = repository.DB.Debug().Create(&account).Error
	} else {
		err = repository.DB.Create(&account).Error
	}

	if err != nil {
		return &models.Account{}, err
	}

	return account, nil
}

func (repository *AccountRepository) Update(account *models.Account, data map[string]interface{}) (*models.Account, error) {
	var err error

	if repository.Debug {
		err = repository.DB.Debug().Model(&account).Where("id = ?", account.ID).Updates(data).Error
	} else {
		err = repository.DB.Model(&account).Where("id = ?", account.ID).Updates(data).Error
	}

	if err != nil {
		return account, err
	}

	if repository.Debug {
		err = repository.DB.Debug().Where("id = ?", account.ID).Take(&account).Error
	} else {
		err = repository.DB.Where("id = ?", account.ID).Take(&account).Error
	}

	if err != nil {
		return account, err
	}

	return account, nil
}

func (repository *AccountRepository) Delete(account *models.Account) (int64, error) {
	var err error
	var db *gorm.DB

	if repository.Debug {
		err = repository.DB.Debug().Model(&account).Association("Roles").Clear()
	} else {
		err = repository.DB.Model(&account).Association("Roles").Clear()
	}

	if err != nil {
		return 0, err
	}

	if repository.Debug {
		db = repository.DB.Debug().Where("id = ?", account.ID).Delete(&models.Account{})
	} else {
		db = repository.DB.Where("id = ?", account.ID).Delete(&models.Account{})
	}

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (repository *AccountRepository) FindAll() (accounts *[]models.Account, err error) {
	if repository.Debug {
		err = repository.DB.Debug().Find(&accounts).Error
	} else {
		err = repository.DB.Find(&accounts).Error
	}

	if err != nil {
		return &[]models.Account{}, err
	}

	return accounts, err
}

func (repository *AccountRepository) FindOneByID(uid uint) (*models.Account, error) {
	var db *gorm.DB
	var account *models.Account

	if repository.Debug {
		db = repository.DB.Debug().Preload("Roles").First(&account, uid)
	} else {
		db = repository.DB.Preload("Roles").First(&account, uid)
	}

	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		return account, fmt.Errorf("account not found: %s", db.Error)
	}

	return account, db.Error
}

func (repository *AccountRepository) FindOneByEmail(email string, withRoles bool) (*models.Account, error) {
	var db *gorm.DB
	var account *models.Account

	if repository.Debug {
		db = repository.DB.Debug().Where("email = ?", email)
	} else {
		db = repository.DB.Where("email = ?", email)
	}

	if withRoles {
		db = db.Preload("Roles").First(&account)
	} else {
		db = db.First(&account)
	}

	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		return account, fmt.Errorf("account not found: %s", db.Error)
	}

	return account, db.Error
}

func (repository *AccountRepository) AddRoles(account *models.Account, roles []models.Role) (*models.Account, error) {
	var err error

	if repository.Debug {
		err = repository.DB.Debug().Model(&account).Association("Roles").Append(&roles)
	} else {
		err = repository.DB.Model(&account).Association("Roles").Append(&roles)
	}

	if err != nil {
		return nil, err
	}

	account, err = repository.FindOneByID(account.ID)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (repository *AccountRepository) DeleteRoles(account *models.Account, roles []models.Role) (*models.Account, error) {
	var err error

	if repository.Debug {
		err = repository.DB.Debug().Model(&account).Association("Roles").Delete(&roles)
	} else {
		err = repository.DB.Model(&account).Association("Roles").Delete(&roles)
	}

	if err != nil {
		return nil, err
	}

	account, err = repository.FindOneByID(account.ID)

	if err != nil {
		return nil, err
	}

	return account, nil
}
