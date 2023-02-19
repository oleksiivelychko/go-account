package models

import (
	"github.com/oleksiivelychko/go-account/utils"
)

type Role struct {
	ID        uint               `gorm:"primary_key;auto_increment" json:"id"`
	Name      string             `gorm:"size:50;not null;unique" json:"name"`
	CreatedAt utils.GormDateTime `json:"created_at"`
	UpdatedAt utils.GormDateTime `json:"updated_at"`
	Accounts  []Account          `gorm:"many2many:accounts_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"accounts,omitempty"`
}

type RoleSerialized struct {
	Name string `json:"name"`
}
