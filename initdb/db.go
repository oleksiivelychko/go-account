package initdb

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"time"
)

func Connection(dsn, logging string) (db *gorm.DB, err error) {
	if logging == "enable" {
		currentTime := time.Now().UTC()
		formatDate := currentTime.Format("01-02-2006")

		var f *os.File
		f, err = os.OpenFile("docker/logs/gorm_"+formatDate+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			f, err = os.OpenFile("../docker/logs/gorm_"+formatDate+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return
			}
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
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			// mute log, including console errors
			Logger: logger.Default.LogMode(logger.Silent),
		})
	}

	return
}
