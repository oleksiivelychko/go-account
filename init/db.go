package init

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Connection(host, user, pass, name, port, ssl, tz, logging string) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, pass, name, port, ssl, tz,
	)

	if logging == "enable" {
		currentTime := time.Now().UTC()
		formatDate := currentTime.Format("01-02-2006")

		wd, _ := os.Getwd()
		parentDirectory := filepath.Dir(wd) + "/"

		var f *os.File
		f, err = os.OpenFile(
			parentDirectory+"docker/logs/gorm_"+formatDate+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return
		}

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			// log to console and file
			Logger: logger.New(
				log.New(io.MultiWriter(os.Stdout, f), "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold:             time.Second,   // Slow SQL threshold
					LogLevel:                  logger.Silent, // Log level
					IgnoreRecordNotFoundError: false,         // Ignore ErrRecordNotFound error for logger
					Colorful:                  false,         // Disable color
				},
			),
		})
	} else {
		db, err = gorm.Open(postgres.Open(dsn))
	}

	/*if cfg.Debug {
		_ = db.Debug().AutoMigrate(
			&models.Account{},
			&models.Role{},
		)
	} else {
		_ = db.AutoMigrate(
			&models.Account{},
			models.Role{},
		)
	}*/

	return
}
