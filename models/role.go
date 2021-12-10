package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type RoleRepository struct {
	DB    *gorm.DB
	Debug bool
}

type Role struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:50;not null;unique" json:"name"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Accounts  []Account `gorm:"many2many:accounts_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type RoleSerialized struct {
	Name string `json:"name"`
}

func (rr *RoleRepository) Validate(role *Role) error {
	if role.Name != "" {
		existsRole := &Role{}
		err := rr.DB.Where("name = ?", role.Name).First(existsRole).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if existsRole.Name != "" {
			return errors.New("role already exists")
		}
	}

	return nil
}

func (rr *RoleRepository) Create(role *Role) (*Role, error) {
	err := rr.Validate(role)
	if err != nil {
		return &Role{}, err
	}

	if rr.Debug {
		err = rr.DB.Debug().Create(&role).Error
	} else {
		err = rr.DB.Create(&role).Error
	}

	if err != nil {
		return &Role{}, err
	}
	return role, nil
}

func (rr *RoleRepository) Update(role *Role) (*Role, error) {
	err := rr.Validate(role)
	if err != nil {
		return role, err
	}

	data := map[string]interface{}{
		"name":       role.Name,
		"updated_at": time.Now(),
	}

	if rr.Debug {
		err = rr.DB.Debug().Model(&role).Where("id = ?", role.ID).Updates(data).Error
	} else {
		err = rr.DB.Model(&role).Where("id = ?", role.ID).Updates(data).Error
	}

	if err != nil {
		return role, err
	}

	if rr.Debug {
		err = rr.DB.Debug().Where("id = ?", role.ID).Take(&role).Error
	} else {
		err = rr.DB.Where("id = ?", role.ID).Take(&role).Error
	}

	if err != nil {
		return role, err
	}
	return role, nil
}

func (rr *RoleRepository) Delete(role *Role) (int64, error) {
	var db *gorm.DB
	if rr.Debug {
		db = rr.DB.Debug().Where("id = ?", role.ID).Delete(&Role{})
	} else {
		db = rr.DB.Where("id = ?", role.ID).Delete(&Role{})
	}

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (rr *RoleRepository) FindAll() (roles *[]Role, err error) {
	if rr.Debug {
		err = rr.DB.Debug().Find(&roles).Error
	} else {
		err = rr.DB.Find(&roles).Error
	}

	if err != nil {
		return &[]Role{}, err
	}
	return roles, err
}

func (rr *RoleRepository) FindOneByID(uid uint) (role *Role, err error) {
	var db *gorm.DB
	if rr.Debug {
		db = rr.DB.Debug().First(&role, uid)
	} else {
		db = rr.DB.First(&role, uid)
	}

	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		return role, fmt.Errorf("role `%d` doesn't managed to find: %s", uid, db.Error)
	}
	return role, db.Error
}

func (rr *RoleRepository) FindOneByName(name string) (role *Role, err error) {
	if rr.Debug {
		err = rr.DB.Debug().Where("name = ?", name).Take(&role).Error
	} else {
		err = rr.DB.Where("name = ?", name).Take(&role).Error
	}

	if err != nil {
		return role, fmt.Errorf("role `%s` doesn't managed to find: %s", name, err)
	}
	return role, nil
}

func (rr *RoleRepository) FindOneByNameOrCreate(name string) (*Role, error) {
	role, err := rr.FindOneByName(name)
	if err != nil {
		role, err = rr.Create(&Role{Name: name})
	}
	return role, nil
}
