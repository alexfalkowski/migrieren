package checker

import (
	"context"

	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// NewNoopChecker is an alias of checker.NewNoopChecker.
func NewNoopChecker() *checker.NoopChecker {
	return checker.NewNoopChecker()
}

// NewMigrator constructs a migration health checker for db.
//
// The checker resolves db source/URL through fs and calls migrator.Ping with
// the provided timeout budget on each check.
func NewMigrator(db *migrate.Database, fs *os.FS, migrator *migrate.Migrator, timeout time.Duration) *Migrator {
	return &Migrator{db: db, fs: fs, migrator: migrator, timeout: timeout}
}

// Migrator validates migration source/database connectivity for one database.
type Migrator struct {
	db       *migrate.Database
	fs       *os.FS
	migrator *migrate.Migrator
	timeout  time.Duration
}

// Check executes the migration connectivity check.
//
// It:
//   - resolves source and URL values from the configured database entry,
//   - applies a timeout to the context,
//   - calls migrator.Ping.
//
// Source/URL read errors and ping errors are returned unchanged.
func (c *Migrator) Check(ctx context.Context) error {
	source, err := c.db.GetSource(c.fs)
	if err != nil {
		return err
	}

	url, err := c.db.GetURL(c.fs)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.migrator.Ping(ctx, bytes.String(source), bytes.String(url))
}
