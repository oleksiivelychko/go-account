package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/oleksiivelychko/go-account/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Session(dsn, dbLogPath string) (*gorm.DB, error) {
	var writerLog io.Writer = os.Stdout
	var logLevel = logger.Silent

	if dbLogPath != "" {
		path, err := GetLogPath(dbLogPath)
		if err != nil {
			return nil, err
		}

		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
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
		dsn = makeDSN(false)
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

func makeDSN(isTest bool) string {
	dsn := "host=%s port=%s dbname=%s user=%s password=%s sslmode=%s TimeZone=%s"

	if isTest {
		return fmt.Sprintf(dsn,
			os.Getenv("TEST_DB_HOST"),
			os.Getenv("TEST_DB_PORT"),
			os.Getenv("TEST_DB_NAME"),
			os.Getenv("TEST_DB_USERNAME"),
			os.Getenv("TEST_DB_PASSWORD"),
			os.Getenv("DB_SSL_MODE"),
			os.Getenv("DB_TIMEZONE"),
		)
	} else {
		return fmt.Sprintf(dsn,
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_SSL_MODE"),
			os.Getenv("DB_TIMEZONE"),
		)
	}
}

func testConnection(dbLogPath string) (*gorm.DB, error) {
	err := godotenv.Load("./../.env.test")
	if err != nil {
		log.Fatal("error loading .env.test file")
	}

	return Session(makeDSN(true), dbLogPath)
}

func PrepareTestDB() (*gorm.DB, error) {
	sessionDB, err := testConnection("")
	err = AutoMigrate(sessionDB)

	statement := "TRUNCATE accounts RESTART IDENTITY CASCADE"
	execSQL := sessionDB.Exec(statement)
	if execSQL.Error != nil {
		return nil, execSQL.Error
	}

	statement = "TRUNCATE roles RESTART IDENTITY CASCADE"
	execSQL = sessionDB.Exec(statement)
	if execSQL.Error != nil {
		return nil, execSQL.Error
	}

	return sessionDB, err
}

func GetLogPath(dbLogPath string) (string, error) {
	absPath, err := filepath.Abs(dbLogPath)
	if err != nil {
		return "", err
	}

	if _, err = os.Stat(absPath); os.IsNotExist(err) {
		if err = os.Mkdir(absPath, os.ModePerm); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%s/gorm_%s.log", absPath, time.Now().UTC().Format("02-01-2006")), nil
}
