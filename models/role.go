package models

import "github.com/oleksiivelychko/go-account/types"

type Role struct {
	ID        uint           `gorm:"primary_key;auto_increment" json:"id"`
	Name      string         `gorm:"size:50;not null;unique" json:"name"`
	CreatedAt types.DateTime `json:"createdAt"`
	UpdatedAt types.DateTime `json:"updatedAt"`
	Accounts  []Account      `gorm:"many2many:accounts_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"accounts,omitempty"`
}

type RoleSerialized struct {
	Name string `json:"name"`
}
