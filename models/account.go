package models

type Account struct {
	ID        uint     `gorm:"primary_key;auto_increment" json:"id"`
	Email     string   `gorm:"size:100;not null;unique" json:"email"`
	Password  string   `gorm:"size:100;not null;" json:"password"`
	CreatedAt JSONTime `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt JSONTime `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Roles     []Role   `gorm:"many2many:accounts_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type AccountSerialized struct {
	ID             uint             `json:"id,omitempty"`
	Email          string           `json:"email,omitempty"`
	Roles          []RoleSerialized `json:"roles,omitempty"`
	AccessToken    string           `json:"access_token,omitempty"`
	RefreshToken   string           `json:"refresh_token,omitempty"`
	ExpirationTime string           `json:"expiration_time,omitempty"`
	ErrorMessage   string           `json:"error_message,omitempty"`
	ErrorCode      uint8            `json:"error_code,omitempty"`
}
