package models

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Role struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:50;not null;unique" json:"name"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Accounts  []Account `gorm:"many2many:accounts_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type RoleRepository struct {
	orm ORM
}

func (rr *RoleRepository) Validate(model *Role) error {
	if model.Name != "" {
		existsRole := &Role{}
		err := rr.orm.DB.Where("name = ?", model.Name).First(existsRole).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if existsRole.Name != "" {
			return errors.New("role already exists")
		}
	}

	return nil
}

func (rr *RoleRepository) Create(model *Role) (*Role, error) {
	err := rr.Validate(model)
	if err != nil {
		return &Role{}, err
	}

	if rr.orm.debug {
		err = rr.orm.DB.Debug().Create(&model).Error
	} else {
		err = rr.orm.DB.Create(&model).Error
	}

	if err != nil {
		return &Role{}, err
	}
	return model, nil
}

func (rr *RoleRepository) Update(model *Role) (*Role, error) {
	err := rr.Validate(model)
	if err != nil {
		return model, err
	}

	data := map[string]interface{}{
		"name":       model.Name,
		"updated_at": time.Now(),
	}

	if rr.orm.debug {
		err = rr.orm.DB.Debug().Model(&model).Where("id = ?", model.ID).Updates(data).Error
	} else {
		err = rr.orm.DB.Model(&model).Where("id = ?", model.ID).Updates(data).Error
	}

	if err != nil {
		return model, err
	}

	if rr.orm.debug {
		err = rr.orm.DB.Debug().Where("id = ?", model.ID).Take(&model).Error
	} else {
		err = rr.orm.DB.Where("id = ?", model.ID).Take(&model).Error
	}

	if err != nil {
		return model, err
	}
	return model, nil
}

func (rr *RoleRepository) Delete(model *Role) (int64, error) {
	var db *gorm.DB
	if rr.orm.debug {
		db = rr.orm.DB.Debug().Where("id = ?", model.ID).Delete(&Role{})
	} else {
		db = rr.orm.DB.Where("id = ?", model.ID).Delete(&Role{})
	}

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (rr *RoleRepository) FindAll() (*[]Role, error) {
	var roles []Role
	var err error
	if rr.orm.debug {
		err = rr.orm.DB.Debug().Find(&roles).Error
	} else {
		err = rr.orm.DB.Find(&roles).Error
	}
	if err != nil {
		return &[]Role{}, err
	}
	return &roles, err
}

func (rr *RoleRepository) FindOneByID(uid uint) (*Role, error) {
	var db *gorm.DB
	var model *Role
	if rr.orm.debug {
		db = rr.orm.DB.Debug().First(&model, uid)
	} else {
		db = rr.orm.DB.First(&model, uid)
	}

	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		return model, errors.New("role not found")
	}
	return model, db.Error
}

func (rr *RoleRepository) FindOneByName(name string) (*Role, error) {
	var model *Role
	var err error

	if rr.orm.debug {
		err = rr.orm.DB.Debug().Where("name = ?", name).Take(&model).Error
	} else {
		err = rr.orm.DB.Debug().Where("name = ?", name).Take(&model).Error
	}

	if err != nil {
		return model, errors.New("role not found")
	}
	return model, nil
}
