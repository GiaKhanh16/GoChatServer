package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// Config holds environment variables
type Config struct {
	DBURL string
}

// LoadConfig loads environment variables and connects to the database
func LoadConfig() (*Config, *pgx.Conn) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: no .env file found, relying on environment variables")
	}

	dbURL := os.Getenv("SUPABASE_DB_URL")
	if dbURL == "" {
		log.Fatal("SUPABASE_DB_URL not set in environment")
	}

	// Connect to Postgres / Supabase
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	cfg := &Config{
		DBURL: dbURL,
	}

	return cfg, conn
}
