package initdb

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Connection(dsn, logging string) (db *gorm.DB, err error) {
	if logging == "enable" {
		currentTime := time.Now().UTC()
		formatDate := currentTime.Format("02-01-2006")

		logPath := filepath.Join(".", "logs")
		_ = os.Mkdir(logPath, 0755)

		var f *os.File
		f, err = os.OpenFile("./logs/gorm_"+formatDate+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			log.Fatalf(err.Error())
		}

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			// To log into console add `os.Stdout` as argument to `io.MultiWriter` function.
			Logger: logger.New(
				log.New(io.MultiWriter(f), "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold:             time.Second,   // Slow SQL threshold
					LogLevel:                  logger.Silent, // Log level
					IgnoreRecordNotFoundError: false,         // Don't ignore ErrRecordNotFound error
					Colorful:                  false,         // Disable color
				},
			),
		})
	} else {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			// mute log, including console errors
			Logger: logger.Default.LogMode(logger.Silent),
		})
	}

	return
}
