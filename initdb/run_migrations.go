package initdb

import (
	"github.com/oleksiivelychko/go-account/models"
	"gorm.io/gorm"
	"os"
)

func AutoMigrate(db *gorm.DB) error {
	if os.Getenv("DB_LOG") == "enable" {
		return db.Debug().AutoMigrate(
			&models.Account{},
			&models.Role{},
		)
	} else {
		return db.AutoMigrate(
			&models.Account{},
			&models.Role{},
		)
	}
}
