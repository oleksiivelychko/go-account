package db

import (
	"fmt"
	"github.com/oleksiivelychko/go-account/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"time"
)

func Session(dsn, dbLogPath string) (*gorm.DB, error) {
	var writerLog io.Writer = os.Stdout
	var logLevel = logger.Silent

	if dbLogPath != "" {
		file, err := os.OpenFile(GetLogPath(dbLogPath), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			return nil, err
		}

		writerLog = io.MultiWriter(file)
		logLevel = logger.Info
	}

	return gorm.Open(postgres.Open(dsn), makeConfig(writerLog, logLevel))
}

func makeConfig(writer io.Writer, level logger.LogLevel) *gorm.Config {
	return &gorm.Config{
		Logger: logger.New(
			log.New(io.MultiWriter(writer), "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second, // slow SQL threshold
				LogLevel:                  level,       // log level
				IgnoreRecordNotFoundError: false,       // do not ignore ErrRecordNotFound error
				Colorful:                  true,        // enable color
			},
		),
	}
}

/*
Connection DATABASE_URL=postgres://admin:secret@localhost:5432/account
*/
func Connection() (*gorm.DB, error) {
	var dsn = os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_SSL_MODE"),
			os.Getenv("DB_TIMEZONE"),
		)
	}

	return Session(dsn, os.Getenv("DB_LOG_PATH"))
}

func AutoMigrate(db *gorm.DB) error {
	if os.Getenv("DB_LOG_PATH") != "" {
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

func makeTestDB(dbLogPath string) (*gorm.DB, error) {
	var dsn = fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s TimeZone=%s",
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_NAME"),
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("DB_SSL_MODE"),
		os.Getenv("DB_TIMEZONE"),
	)

	return Session(dsn, dbLogPath)
}

func PrepareTestDB() (*gorm.DB, error) {
	db, err := makeTestDB("")
	err = AutoMigrate(db)

	statement := "TRUNCATE accounts RESTART IDENTITY CASCADE"
	execSQL := db.Exec(statement)
	if execSQL.Error != nil {
		return nil, execSQL.Error
	}

	statement = "TRUNCATE roles RESTART IDENTITY CASCADE"
	execSQL = db.Exec(statement)
	if execSQL.Error != nil {
		return nil, execSQL.Error
	}

	return db, err
}

func GetLogPath(basePath string) string {
	return fmt.Sprintf("./%s/gorm_%s.log", basePath, time.Now().UTC().Format("02-01-2006"))
}
