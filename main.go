package main

import (
	"database/sql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/oleksiivelychko/go-account/handlers"
	"github.com/oleksiivelychko/go-account/initdb"
	"log"
	"net/http"
	"os"
)

func main() {
	// initdb.LoadEnv() // uncomment to local development

	db, err := initdb.DB()
	if err != nil {
		log.Fatalf("failed database connection: %s", err)
	}

	dbConnection, err := db.DB()
	defer func(sqlDB *sql.DB) {
		err = sqlDB.Close()
		if err != nil {
			log.Fatalf("unable to close database connection: %s", err)
		}
	}(dbConnection)

	http.HandleFunc("/api/account/register", handlers.RegisterHandler(db))
	http.HandleFunc("/api/account/login", handlers.LoginHandler(db))
	http.HandleFunc("/api/account/user", handlers.UserHandler(db))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
