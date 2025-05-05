package migrate

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sqlx.DB, migrationsDir string) error {
	goose.SetLogger(goose.NopLogger()) // Or use log.Default() if you want logging
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}
	if err := goose.Up(db.DB, migrationsDir); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}
