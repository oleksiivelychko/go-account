package initdb

import (
	"fmt"
	"gorm.io/gorm"
	"os"
)

func LoadEnv() {
	if host, isSet := os.LookupEnv("HOST"); isSet == true {
		_ = os.Setenv("HOST", host)
	} else {
		_ = os.Setenv("HOST", "localhost")
	}
	if os.Getenv("PORT") == "" {
		_ = os.Setenv("PORT", "8081")
	}
	if os.Getenv("DB_LOG") == "" {
		_ = os.Setenv("DB_LOG", "disable")
	}
	if host, isSet := os.LookupEnv("DB_HOST"); isSet == true {
		_ = os.Setenv("DB_HOST", host)
	} else {
		_ = os.Setenv("DB_HOST", "localhost")
	}
	if os.Getenv("DB_PORT") == "" {
		_ = os.Setenv("DB_PORT", "5432")
	}
	if os.Getenv("DB_NAME") == "" {
		_ = os.Setenv("DB_NAME", "account")
	}
	if os.Getenv("DB_USER") == "" {
		_ = os.Setenv("DB_USER", "admin")
	}
	if os.Getenv("DB_PASS") == "" {
		_ = os.Setenv("DB_PASS", "secret")
	}
	if os.Getenv("DB_DRIVER") == "" {
		_ = os.Setenv("DB_DRIVER", "postgres")
	}
	if os.Getenv("DB_SSL") == "" {
		_ = os.Setenv("DB_SSL", "disable")
	}
	if os.Getenv("DB_TZ") == "" {
		_ = os.Setenv("DB_TZ", "UTC")
	}
	if os.Getenv("DATABASE_URL") == "" {
		_ = os.Setenv("DATABASE_URL", "postgres://gopher:secret@localhost:5432/go-postgres")
	}
	if os.Getenv("TEST_DB_HOST") == "" {
		_ = os.Setenv("TEST_DB_HOST", "localhost")
	}
	if os.Getenv("TEST_DB_PORT") == "" {
		_ = os.Setenv("TEST_DB_PORT", "5433")
	}
	if os.Getenv("TEST_DB_NAME") == "" {
		_ = os.Setenv("TEST_DB_NAME", "account-test")
	}
	if os.Getenv("TEST_DB_USER") == "" {
		_ = os.Setenv("TEST_DB_USER", "test")
	}
	if os.Getenv("TEST_DB_PASS") == "" {
		_ = os.Setenv("TEST_DB_PASS", "secret")
	}
	if os.Getenv("APP_JWT_URL") == "" {
		_ = os.Setenv("APP_JWT_URL", "http://localhost:8080")
	}
}

func DB() (*gorm.DB, error) {
	var dsn = os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_SSL"),
			os.Getenv("DB_TZ"),
		)
	}

	return Connection(dsn, os.Getenv("DB_LOG"))
}

func TestDB() (*gorm.DB, error) {
	var dsn = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASS"),
		os.Getenv("TEST_DB_NAME"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("DB_SSL"),
		os.Getenv("DB_TZ"),
	)
	return Connection(dsn, os.Getenv("DB_LOG"))
}
