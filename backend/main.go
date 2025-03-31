package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/pdridh/service-needs-app/backend/auth"
	"github.com/pdridh/service-needs-app/backend/config"
	"github.com/pdridh/service-needs-app/backend/db"
	"github.com/pdridh/service-needs-app/backend/server"
	"github.com/pdridh/service-needs-app/backend/user"
)

func main() {
	config.Load()

	// Connect to db
	db.ConnectToDB()
	defer db.DisconnectFromDB()

	userStore := user.NewMongoStore(db.GetCollectionFromDB(config.Server().DatabaseName, config.Server().UserCollectionName))
	validate := validator.New()
	authService := auth.NewService(userStore, validate)
	authHandler := auth.NewHandler(authService)

	srv := server.New(authHandler)

	httpServer := &http.Server{
		Addr:         net.JoinHostPort(config.Server().Host, config.Server().Port),
		Handler:      srv,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listening on %s\n", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}
}
