package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pdridh/service-needs-app/backend/auth"
	"github.com/pdridh/service-needs-app/backend/business"
	"github.com/pdridh/service-needs-app/backend/config"
	"github.com/pdridh/service-needs-app/backend/consumer"
	"github.com/pdridh/service-needs-app/backend/db"
	"github.com/pdridh/service-needs-app/backend/server"
	"github.com/pdridh/service-needs-app/backend/user"
)

func main() {
	config.Load()

	// Connect to db
	db.ConnectToDB()
	defer db.DisconnectFromDB()

	validate := validator.New()

	userStore := user.NewMongoStore(db.GetCollectionFromDB(config.Server().DatabaseName, config.Server().UserCollectionName))
	businessStore := business.NewMongoStore(db.GetCollectionFromDB(config.Server().DatabaseName, config.Server().BusinessCollectionName))
	consumerStore := consumer.NewMongoStore(db.GetCollectionFromDB(config.Server().DatabaseName, config.Server().ConsumerCollectionName))
	businessService := business.NewService(businessStore, validate)
	businessHandler := business.NewHandler(businessService)

	authService := auth.NewService(db.GetClient(), userStore, businessStore, consumerStore, validate)
	authHandler := auth.NewHandler(authService)

	srv := server.New(authHandler, businessHandler)

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
