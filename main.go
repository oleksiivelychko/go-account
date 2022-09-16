package main

import (
	"context"
	"database/sql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/oleksiivelychko/go-account/handlers"
	"github.com/oleksiivelychko/go-account/initdb"
	"github.com/oleksiivelychko/go-account/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	initdb.LoadEnv()
	addr := os.Getenv("HOST") + ":" + os.Getenv("PORT")

	db, err := initdb.DB()
	if err != nil {
		log.Fatalf("Failed database connection: %s", err)
	}

	err = models.AutoMigrate(db)
	if err != nil {
		log.Fatalf("Failed to migrate models: %s", err)
	}

	dbConnection, err := db.DB()
	defer func(sqlDB *sql.DB) {
		err = sqlDB.Close()
		if err != nil {
			log.Fatalf("Unable to close database connection: %s", err)
		}
	}(dbConnection)

	serveMux := http.NewServeMux()
	serveMux.Handle("/api/account/register/", handlers.NewRegisterHandler(db))
	serveMux.Handle("/api/account/login/", handlers.NewLoginHandler(db))
	serveMux.Handle("/api/account/user/", handlers.NewUserHandler(db))

	server := &http.Server{
		Addr:         addr,
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		log.Printf("Starting server on %s", addr)
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
	_ = server.Shutdown(ctx)
}
