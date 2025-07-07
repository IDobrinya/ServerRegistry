package main

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
	"server-registry/internal/database"
	"time"

	"server-registry/internal/config"
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
			logger.Fatalf("Failed to close database connection: %v", err)
		} else {
			logger.Println("Database connection closed successfully.")
		}
	}(db)

	logger.Println("Database connected successfully")
	logger.Printf("Server will start on port %d", cfg.Port)
}
