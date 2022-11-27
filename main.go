package main

import (
	"context"
	"database/sql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/oleksiivelychko/go-account/handlers"
	"github.com/oleksiivelychko/go-account/initdb"
	"github.com/oleksiivelychko/go-account/repositories"
	"github.com/oleksiivelychko/go-account/services"
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
		log.Fatalf("failed database connection: %s", err)
	}

	err = initdb.AutoMigrate(db)
	if err != nil {
		log.Fatalf("failed to run migrations: %s", err)
	}

	dbConnection, err := db.DB()
	defer func(sqlDB *sql.DB) {
		err = sqlDB.Close()
		if err != nil {
			log.Fatalf("unable to close database connection: %s", err)
		}
	}(dbConnection)

	repository := repositories.NewRepository(db, true)
	accountRepository := repositories.NewAccountRepository(repository)
	accountService := services.NewAccountService(accountRepository)
	roleRepository := repositories.NewRoleRepository(repository)
	roleService := services.NewRoleService(roleRepository)

	serveMux := http.NewServeMux()
	serveMux.Handle("/api/account/register/", handlers.NewRegisterHandler(accountService, roleService))
	serveMux.Handle("/api/account/login/", handlers.NewLoginHandler(accountService))
	serveMux.Handle("/api/account/user/", handlers.NewUserHandler(accountService))

	server := &http.Server{
		Addr:         addr,
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("starting server on %s", addr)
		err = server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	log.Println("received terminate, graceful shutdown", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_ = server.Shutdown(ctx)
}
