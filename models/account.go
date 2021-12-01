package models

import (
	"errors"
	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type AccountRepository struct {
	DB    *gorm.DB
	Debug bool
}

type Account struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Roles     []Role    `gorm:"many2many:accounts_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (model *Account) VerifyPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(model.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (ar *AccountRepository) Validate(account *Account) error {
	if account.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		account.Password = string(hashedPassword)
	}

	if account.Email != "" {
		if err := checkmail.ValidateFormat(account.Email); err != nil {
			return errors.New("invalid:email")
		}

		existsAccount := &Account{}
		err := ar.DB.Where("email = ?", account.Email).First(existsAccount).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if existsAccount.Email != "" {
			return errors.New("email address already exists")
		}
	}

	return nil
}

func (ar *AccountRepository) Create(account *Account) (*Account, error) {
	err := ar.Validate(account)
	if err != nil {
		return &Account{}, err
	}

	if ar.Debug {
		err = ar.DB.Debug().Create(&account).Error
	} else {
		err = ar.DB.Create(&account).Error
	}

	if err != nil {
		return &Account{}, err
	}
	return account, nil
}

func (ar *AccountRepository) Update(account *Account) (*Account, error) {
	err := ar.Validate(account)
	if err != nil {
		return account, err
	}

	data := map[string]interface{}{
		"password":   account.Password,
		"email":      account.Email,
		"updated_at": time.Now(),
	}

	if ar.Debug {
		err = ar.DB.Debug().Model(&account).Where("id = ?", account.ID).Updates(data).Error
	} else {
		err = ar.DB.Model(&account).Where("id = ?", account.ID).Updates(data).Error
	}

	if err != nil {
		return account, err
	}

	if ar.Debug {
		err = ar.DB.Debug().Where("id = ?", account.ID).Take(&account).Error
	} else {
		err = ar.DB.Where("id = ?", account.ID).Take(&account).Error
	}

	if err != nil {
		return account, err
	}
	return account, nil
}

func (ar *AccountRepository) Delete(account *Account) (int64, error) {
	var err error
	var db *gorm.DB

	if ar.Debug {
		err = ar.DB.Debug().Model(&account).Association("Roles").Clear()
	} else {
		err = ar.DB.Model(&account).Association("Roles").Clear()
	}

	if err != nil {
		return 0, err
	}

	if ar.Debug {
		db = ar.DB.Debug().Where("id = ?", account.ID).Delete(&Account{})
	} else {
		db = ar.DB.Where("id = ?", account.ID).Delete(&Account{})
	}

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (ar *AccountRepository) FindAll() (accounts *[]Account, err error) {
	if ar.Debug {
		err = ar.DB.Debug().Find(&accounts).Error
	} else {
		err = ar.DB.Find(&accounts).Error
	}

	if err != nil {
		return &[]Account{}, err
	}
	return accounts, err
}

func (ar *AccountRepository) FindOneByID(uid uint) (*Account, error) {
	var db *gorm.DB
	var account *Account

	if ar.Debug {
		db = ar.DB.Debug().Preload("Roles").First(&account, uid)
	} else {
		db = ar.DB.Preload("Roles").First(&account, uid)
	}

	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		return account, errors.New("account not found")
	}
	return account, db.Error
}

func (ar *AccountRepository) AddRoles(account *Account, roles []Role) (*Account, error) {
	var err error
	if ar.Debug {
		err = ar.DB.Debug().Model(&account).Association("Roles").Append(&roles)
	} else {
		err = ar.DB.Model(&account).Association("Roles").Append(&roles)
	}

	if err != nil {
		return nil, err
	}

	account, err = ar.FindOneByID(account.ID)

	if err != nil {
		return nil, err
	}
	return account, nil
}

func (ar *AccountRepository) DeleteRoles(account *Account, roles []Role) (*Account, error) {
	var err error
	if ar.Debug {
		err = ar.DB.Debug().Model(&account).Association("Roles").Delete(&roles)
	} else {
		err = ar.DB.Model(&account).Association("Roles").Delete(&roles)
	}

	if err != nil {
		return nil, err
	}

	account, err = ar.FindOneByID(account.ID)

	if err != nil {
		return nil, err
	}
	return account, nil
}
