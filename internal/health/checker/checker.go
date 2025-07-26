package checker

import (
	"context"

	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/migrieren/internal/migrate"
	"github.com/alexfalkowski/migrieren/internal/migrate/migrator"
)

// NewNoopChecker is an alias of checker.NewNoopChecker.
func NewNoopChecker() *checker.NoopChecker {
	return checker.NewNoopChecker()
}

// NewMigrator checker.
func NewMigrator(db *migrate.Database, fs *os.FS, migrator migrator.Migrator) *Migrator {
	return &Migrator{db: db, fs: fs, migrator: migrator}
}

// Migrator checker.
type Migrator struct {
	db       *migrate.Database
	fs       *os.FS
	migrator migrator.Migrator
}

// Check the migrator.
func (c *Migrator) Check(ctx context.Context) error {
	source, err := c.db.GetSource(c.fs)
	if err != nil {
		return err
	}

	url, err := c.db.GetURL(c.fs)
	if err != nil {
		return err
	}

	return c.migrator.Ping(ctx, bytes.String(source), bytes.String(url))
}
