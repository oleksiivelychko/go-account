package models

import (
	"github.com/oleksiivelychko/go-account/utils"
)

type Account struct {
	ID        uint           `gorm:"primary_key;auto_increment" json:"id"`
	Email     string         `gorm:"size:100;not null;unique" json:"email"`
	Password  string         `gorm:"size:100;not null;" json:"password"`
	CreatedAt utils.DateTime `json:"createdAt"`
	UpdatedAt utils.DateTime `json:"updatedAt"`
	Roles     []Role         `gorm:"many2many:accounts_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type AccountSerialized struct {
	ID             uint             `json:"id,omitempty"`
	Email          string           `json:"email,omitempty"`
	Roles          []RoleSerialized `json:"roles,omitempty"`
	AccessToken    string           `json:"accessToken,omitempty"`
	RefreshToken   string           `json:"refreshToken,omitempty"`
	ExpirationTime string           `json:"expirationTime,omitempty"`
	ErrorMessage   string           `json:"errorMessage,omitempty"`
}
