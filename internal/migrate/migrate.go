package migrate

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB, migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}
	return goose.Up(db, migrationsDir)
}
