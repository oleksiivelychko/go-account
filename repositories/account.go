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

func (repo *Account) Create(account *models.Account) (err error) {
	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Create(&account).Error
	} else {
		err = repo.Repo.DB.Create(&account).Error
	}

	return err
}

func (repo *Account) Update(account *models.Account, data map[string]interface{}) (err error) {
	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Model(&account).Where("id = ?", account.ID).Updates(data).Error
	} else {
		err = repo.Repo.DB.Model(&account).Where("id = ?", account.ID).Updates(data).Error
	}

	if err != nil {
		return err
	}

	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Where("id = ?", account.ID).Take(&account).Error
	} else {
		err = repo.Repo.DB.Where("id = ?", account.ID).Take(&account).Error
	}

	return err
}

func (repo *Account) Delete(account *models.Account) (int64, error) {
	var err error = nil
	var db *gorm.DB

	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Model(&account).Association("Roles").Clear()
	} else {
		err = repo.Repo.DB.Model(&account).Association("Roles").Clear()
	}

	if err != nil {
		return 0, err
	}

	if repo.Repo.Debug {
		db = repo.Repo.DB.Debug().Where("id = ?", account.ID).Delete(&models.Account{})
	} else {
		db = repo.Repo.DB.Where("id = ?", account.ID).Delete(&models.Account{})
	}

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (repo *Account) FindAll() (account *[]models.Account, err error) {
	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Preload("Roles").Find(&account).Error
	} else {
		err = repo.Repo.DB.Preload("Roles").Find(&account).Preload("Roles").Error
	}

	return account, err
}

func (repo *Account) FindOneByID(id uint) (account *models.Account, err error) {
	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Preload("Roles").First(&account, id).Error
	} else {
		err = repo.Repo.DB.Preload("Roles").First(&account, id).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("unable to find account by ID %d: %s", id, err)
	}

	return account, err
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
		return nil, db.Error
	}

	return modelAccount, db.Error
}

func (repo *Account) AddRoles(account *models.Account, roles []*models.Role) (*models.Account, error) {
	var err error

	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Model(&account).Association("Roles").Append(&roles)
	} else {
		err = repo.Repo.DB.Model(&account).Association("Roles").Append(&roles)
	}

	if err != nil {
		return nil, err
	}

	return repo.FindOneByID(account.ID)
}

func (repo *Account) DeleteRoles(account *models.Account, roles []models.Role) (*models.Account, error) {
	var err error

	if repo.Repo.Debug {
		err = repo.Repo.DB.Debug().Model(&account).Association("Roles").Delete(&roles)
	} else {
		err = repo.Repo.DB.Model(&account).Association("Roles").Delete(&roles)
	}

	if err != nil {
		return nil, err
	}

	return repo.FindOneByID(account.ID)
}

func (repo *Account) HasRoles(account *models.Account, roles []string) bool {
	var count int64

	if repo.Repo.Debug {
		count = repo.Repo.DB.Debug().Model(&account).Where("name IN ?", roles).Association("Roles").Count()
	} else {
		count = repo.Repo.DB.Model(&account).Where("name IN ?", roles).Association("Roles").Count()
	}

	return count == int64(len(roles))
}
