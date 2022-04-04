package main

import (
	"context"
	"database/sql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/oleksiivelychko/go-account/handlers"
	"github.com/oleksiivelychko/go-account/initdb"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-helper/env"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	initdb.LoadEnv()

	db, err := initdb.DB()
	if err != nil {
		log.Fatalf("failed database connection: %s", err)
	}

	err = models.AutoMigrate(db)
	if err != nil {
		log.Fatalf("failed to migrate models: %s", err)
	}

	dbConnection, err := db.DB()
	defer func(sqlDB *sql.DB) {
		err = sqlDB.Close()
		if err != nil {
			log.Fatalf("unable to close database connection: %s", err)
		}
	}(dbConnection)

	serveMux := http.NewServeMux()
	serveMux.Handle("/api/account/register/", handlers.NewRegisterHandler(db))
	serveMux.Handle("/api/account/login/", handlers.NewLoginHandler(db))
	serveMux.Handle("/api/account/user/", handlers.NewUserHandler(db))

	server := &http.Server{
		Addr:         env.GetAddr(),
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		log.Printf("Starting server on %s", env.GetAddr())
		err = server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	log.Println("Received terminate, graceful shutdown", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
