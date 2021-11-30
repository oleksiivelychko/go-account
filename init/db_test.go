package init

import (
	"testing"
)

func TestDbConnection(t *testing.T) {
	Env()
	db, err := TestDB()

	if err != nil {
		t.Errorf("unable to init db connection: %s", err)
	}

	_, err = db.DB()
	if err != nil {
		t.Errorf("unable to get db instance: %s", err)
	}
}
