package db

import (
	"errors"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sqlx.DB

func Connect() error {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, relying on system env")
	}

	// Read DSN from environment
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return errors.New("DATABASE_URL is not set")
	}

	// Connect to PostgreSQL using pgx driver
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return err
	}

	// Assign to global variable
	DB = db
	return nil
}
