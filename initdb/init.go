package initdb

import (
	"gorm.io/gorm"
	"os"
)

func LoadEnv() {
	_ = os.Setenv("PORT", "8081")
	_ = os.Setenv("DB_LOG", "disable")
	_ = os.Setenv("DB_HOST", "localhost")
	_ = os.Setenv("DB_PORT", "5432")
	_ = os.Setenv("DB_NAME", "go-postgres")
	_ = os.Setenv("DB_USER", "gopher")
	_ = os.Setenv("DB_PASS", "secret")
	_ = os.Setenv("DB_DRIVER", "postgres")
	_ = os.Setenv("DB_SSL", "disable")
	_ = os.Setenv("DB_TZ", "UTC")
	_ = os.Setenv("TEST_DB_HOST", "localhost")
	_ = os.Setenv("TEST_DB_PORT", "5433")
	_ = os.Setenv("TEST_DB_NAME", "go-postgres-test")
	_ = os.Setenv("TEST_DB_USER", "gopher")
	_ = os.Setenv("TEST_DB_PASS", "secret")
}

func DB() (*gorm.DB, error) {
	return Connection(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL"),
		os.Getenv("DB_TZ"),
		os.Getenv("DB_LOG"),
	)
}

func TestDB() (*gorm.DB, error) {
	return Connection(
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASS"),
		os.Getenv("TEST_DB_NAME"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("DB_SSL"),
		os.Getenv("DB_TZ"),
		os.Getenv("DB_LOG"),
	)
}
