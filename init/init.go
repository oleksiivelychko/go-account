package init

import "os"

func Env() {
	_ = os.Setenv("PORT", "8081")
	_ = os.Setenv("DB_LOG", "enable")
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
