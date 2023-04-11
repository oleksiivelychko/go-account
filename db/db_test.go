package db

import (
	"errors"
	"os"
	"testing"
)

func TestDB_Session(t *testing.T) {
	session, err := makeTestDB("")
	if err != nil {
		t.Fatal(err)
	}

	_, err = session.DB()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDB_SessionLog(t *testing.T) {
	dbLogPath := "./../.data/log"

	_, err := makeTestDB(dbLogPath)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = os.Stat(GetLogPath(dbLogPath)); errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
}

func TestDB_AutoMigrate(t *testing.T) {
	session, err := makeTestDB("")
	if err != nil {
		t.Fatal(err)
	}

	err = AutoMigrate(session)
	if err != nil {
		t.Fatal(err)
	}
}