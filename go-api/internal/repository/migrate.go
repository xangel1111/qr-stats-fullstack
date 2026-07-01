package repository

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres driver
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Migrate applies all pending "up" migrations from the given filesystem
// (the embedded db/migrations directory). It is idempotent: running it when
// the schema is already current is a no-op.
func Migrate(databaseURL string, migrations fs.FS) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("open migrations source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, databaseURL)
	if err != nil {
		return fmt.Errorf("init migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
