package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/jmoiron/sqlx"
	"server-registry/internal/api"
	"server-registry/internal/config"
	"server-registry/internal/database"
)

func main() {
	ctx := context.Background()

	logger := log.New(log.Writer(), "[SERVER] ", log.LstdFlags)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	logger.Println("Initializing database connection")
	db, err := database.NewPostgresDB(ctx, database.PostgresConfig{
		DSN:                cfg.DatabaseURL,
		MaxOpenConnections: 25,
		MaxIdleConnections: 25,
		ConnMaxLifetime:    5 * time.Minute,
	})
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			logger.Printf("Failed to close database connection: %v", err)
		} else {
			logger.Println("Database connection closed successfully.")
		}
	}(db)

	logger.Println("Database connected successfully")

	router := api.NewRouter(db)

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "User-ID"}),
	)(router)

	addr := fmt.Sprintf(":%d", cfg.Port)
	logger.Printf("Starting server on %s", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: corsHandler,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
