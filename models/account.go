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

func (ar *AccountRepository) Validate(model *Account) error {
	if model.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		model.Password = string(hashedPassword)
	}

	if model.Email != "" {
		if err := checkmail.ValidateFormat(model.Email); err != nil {
			return errors.New("invalid:email")
		}

		existsAccount := &Account{}
		err := ar.DB.Where("email = ?", model.Email).First(existsAccount).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if existsAccount.Email != "" {
			return errors.New("email address already exists")
		}
	}

	return nil
}

func (ar *AccountRepository) Create(model *Account) (*Account, error) {
	err := ar.Validate(model)
	if err != nil {
		return &Account{}, err
	}

	if ar.Debug {
		err = ar.DB.Debug().Create(&model).Error
	} else {
		err = ar.DB.Create(&model).Error
	}

	if err != nil {
		return &Account{}, err
	}
	return model, nil
}

func (ar *AccountRepository) Update(model *Account) (*Account, error) {
	err := ar.Validate(model)
	if err != nil {
		return model, err
	}

	data := map[string]interface{}{
		"password":   model.Password,
		"email":      model.Email,
		"updated_at": time.Now(),
	}

	if ar.Debug {
		err = ar.DB.Debug().Model(&model).Where("id = ?", model.ID).Updates(data).Error
	} else {
		err = ar.DB.Model(&model).Where("id = ?", model.ID).Updates(data).Error
	}

	if err != nil {
		return model, err
	}

	if ar.Debug {
		err = ar.DB.Debug().Where("id = ?", model.ID).Take(&model).Error
	} else {
		err = ar.DB.Where("id = ?", model.ID).Take(&model).Error
	}

	if err != nil {
		return model, err
	}
	return model, nil
}

func (ar *AccountRepository) Delete(model *Account) (int64, error) {
	var err error
	var db *gorm.DB

	if ar.Debug {
		err = ar.DB.Debug().Model(&model).Association("Roles").Clear()
	} else {
		err = ar.DB.Model(&model).Association("Roles").Clear()
	}

	if err != nil {
		return 0, err
	}

	if ar.Debug {
		db = ar.DB.Debug().Where("id = ?", model.ID).Delete(&Account{})
	} else {
		db = ar.DB.Where("id = ?", model.ID).Delete(&Account{})
	}

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (ar *AccountRepository) FindAll() (*[]Account, error) {
	var accounts []Account
	var err error
	if ar.Debug {
		err = ar.DB.Debug().Find(&accounts).Error
	} else {
		err = ar.DB.Find(&accounts).Error
	}
	if err != nil {
		return &[]Account{}, err
	}
	return &accounts, err
}

func (ar *AccountRepository) FindOneByID(uid uint) (*Account, error) {
	var db *gorm.DB
	var model *Account

	if ar.Debug {
		db = ar.DB.Debug().Preload("Roles").First(&model, uid)
	} else {
		db = ar.DB.Preload("Roles").First(&model, uid)
	}

	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		return model, errors.New("account not found")
	}
	return model, db.Error
}

func (ar *AccountRepository) AddRoles(model *Account, roles []Role) (*Account, error) {
	var err error
	if ar.Debug {
		err = ar.DB.Debug().Model(&model).Association("Roles").Append(&roles)
	} else {
		err = ar.DB.Model(&model).Association("Roles").Append(&roles)
	}

	if err != nil {
		return nil, err
	}

	model, err = ar.FindOneByID(model.ID)

	if err != nil {
		return nil, err
	}

	return model, nil
}

func (ar *AccountRepository) DeleteRoles(model *Account, roles []Role) (*Account, error) {
	var err error
	if ar.Debug {
		err = ar.DB.Debug().Model(&model).Association("Roles").Delete(&roles)
	} else {
		err = ar.DB.Model(&model).Association("Roles").Delete(&roles)
	}

	if err != nil {
		return nil, err
	}

	model, err = ar.FindOneByID(model.ID)

	if err != nil {
		return nil, err
	}

	return model, nil
}
