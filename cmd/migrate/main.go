package main

import (
	"log"

	"github.com/Alexx1088/userservice/internal/db"
	"github.com/Alexx1088/userservice/internal/migrate"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer db.Pool.Close()

	migrationsPath := "file://migrations"

	if err := migrate.RunMigrations(db.Pool, migrationsPath); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("✅ Migrations applied successfully")
}
