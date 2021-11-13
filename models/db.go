package models

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

type ORM struct {
	DB    *gorm.DB
	debug bool
}

func InitDB(test bool) *ORM {
	var dsn string
	if test {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			os.Getenv("TEST_DB_HOST"),
			os.Getenv("TEST_DB_USER"),
			os.Getenv("TEST_DB_PASS"),
			os.Getenv("TEST_DB_NAME"),
			os.Getenv("TEST_DB_PORT"),
			os.Getenv("DB_SSL"),
			os.Getenv("DB_TIMEZONE"),
		)
	} else {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_SSL"),
			os.Getenv("DB_TIMEZONE"),
		)
	}

	var db *gorm.DB = nil
	var err error
	var debugMode = false

	if os.Getenv("DB_DEBUG") == "enabled" {
		debugMode = true
		currentTime := time.Now().UTC()
		formatDate := currentTime.Format("01-02-2006")

		var parentDirectory = ""
		if test {
			wd, _ := os.Getwd()
			parentDirectory = filepath.Dir(wd) + "/"
		}

		f, err := os.OpenFile(
			parentDirectory+"docker/logs/gorm_"+formatDate+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
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

	var connectedDb = func(testMode bool) string {
		if testMode {
			return os.Getenv("TEST_DB_NAME")
		} else {
			return os.Getenv("DB_NAME")
		}
	}(test)

	if err != nil {
		fmt.Printf("cannot connect to `%s` database\n", connectedDb)
		log.Fatal("database error:", err)
	} else {
		fmt.Printf("connected to the `%s` database", connectedDb)
	}

	if os.Getenv("DB_DEBUG") == "enabled" {
		_ = db.Debug().AutoMigrate(
			&Account{},
			&Role{},
		)
	} else {
		_ = db.AutoMigrate(
			&Account{},
			&Role{},
		)
	}

	return &ORM{db, debugMode}
}
