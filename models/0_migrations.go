package models

import (
	"gorm.io/gorm"
	"os"
)

func AutoMigrate(db *gorm.DB) error {
	if os.Getenv("DB_LOG") == "enable" {
		return db.Debug().AutoMigrate(
			&Account{},
			&Role{},
		)
	} else {
		return db.AutoMigrate(
			&Account{},
			&Role{},
		)
	}
}
