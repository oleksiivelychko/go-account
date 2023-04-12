package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/oleksiivelychko/go-account/db"
	"github.com/oleksiivelychko/go-account/handlers"
	"github.com/oleksiivelychko/go-account/repositories"
	"github.com/oleksiivelychko/go-account/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	session, err := db.Connection()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(session)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := session.DB()
	defer func(conn *sql.DB) {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	repo := repositories.NewRepository(session, true)
	accountRepo := repositories.NewAccount(repo)
	accountService := services.NewAccount(accountRepo)
	roleRepo := repositories.NewRole(repo)
	roleService := services.NewRole(roleRepo)

	serveMux := http.NewServeMux()
	serveMux.Handle("/api/account/register/", handlers.NewRegister(accountService, roleService))
	serveMux.Handle("/api/account/login/", handlers.NewLogin(accountService))
	serveMux.Handle("/api/account/user/", handlers.NewUser(accountService))

	addr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
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

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt)
	signal.Notify(signalCh, os.Kill)

	sig := <-signalCh
	log.Println("received terminate, graceful shutdown", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_ = server.Shutdown(ctx)
}
