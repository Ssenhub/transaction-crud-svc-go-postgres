package main

import (
	"fmt"
	"net/http"
	"os"
	"transaction-crud-svc-go-postgres/handlers"
	"transaction-crud-svc-go-postgres/middlewares"
	"transaction-crud-svc-go-postgres/models"
	"transaction-crud-svc-go-postgres/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("Failed to load .env", err)
		return
	}

	config := &storage.Config{
		DockerHost:   os.Getenv("POSTGRES_DOCKER_HOST"),
		InternalHost: os.Getenv("POSTGRES_INTERNAL_HOST"),
		LocalHost:    os.Getenv("POSTGRES_LOCAL_HOST"),
		Port:         os.Getenv("POSTGRES_PORT"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		User:         os.Getenv("POSTGRES_USER"),
		SslMode:      os.Getenv("POSTGRES_SSLMODE"),
		DbName:       os.Getenv("POSTGRES_DB"),
	}

	db, dbErr := storage.NewConnection(config)
	if dbErr != nil {
		fmt.Println("Failed to connect to DB", dbErr)
		return
	}

	dbErr = models.MigrateTransaction(db)
	if dbErr != nil {
		fmt.Println("MigrateTransaction failed", dbErr)
		return
	}

	var store = storage.PostgresTransactionStore{
		DB: db,
	}

	handlers.DataHandler = handlers.DataStore{
		Store: &store,
	}

	router := chi.NewRouter()

	router.Post("/login", handlers.LoginHandler)

	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middlewares.JwtAutheticator)
		r.Get("/transactions", handlers.GetListHandler)
		r.Get("/transactions/{id}", handlers.GetHandler)
		r.Post("/transactions", handlers.PostHandler)
		r.Put("/transactions/{id}", handlers.PutHandler)
		r.Delete("/transactions/{id}", handlers.DeleteHandler)
		r.Delete("/transactions", handlers.DeleteAllHandler)
	})

	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	listenerErr := server.ListenAndServe()
	if listenerErr != nil {
		fmt.Println("Failed to listen to server", listenerErr)
		return
	}
}
