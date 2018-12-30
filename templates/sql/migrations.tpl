package sql

import (
	"errors"
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/{{ .Driver }}"

	// file driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var Migrations = "migrations"

// RunMigrations performs any required database migrations
func RunMigrations() error {
	if db == nil {
		return errors.New("database not initialized")
	}

	if err := db.Ping(); err != nil {
		return err
	}

	migrationPath := fmt.Sprintf("file://%s", Migrations)
	driver, err := {{ .Driver }}.WithInstance(db, &{{ .Driver }}.Config{})
	m, err := migrate.NewWithDatabaseInstance(migrationPath, dbType, driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return err
	}
	return nil
}