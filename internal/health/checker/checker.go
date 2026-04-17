package checker

import (
	"context"

	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// NewNoopChecker returns a checker that always reports healthy.
func NewNoopChecker() *checker.NoopChecker {
	return checker.NewNoopChecker()
}

// NewMigrator constructs a health checker for one configured migration target.
//
// The checker resolves the database source and URL through fs, then pings the
// migrator with the provided timeout.
func NewMigrator(db *migrate.Database, fs *os.FS, migrator *migrate.Migrator, timeout time.Duration) *Migrator {
	return &Migrator{db: db, fs: fs, migrator: migrator, timeout: timeout}
}

// Migrator checks whether a configured migration target can be resolved and
// pinged successfully.
type Migrator struct {
	db       *migrate.Database
	fs       *os.FS
	migrator *migrate.Migrator
	timeout  time.Duration
}

// Check resolves the migration source and database URL, then pings the target
// database within the configured timeout.
func (c *Migrator) Check(ctx context.Context) error {
	source, err := c.db.GetSource(c.fs)
	if err != nil {
		return err
	}

	url, err := c.db.GetURL(c.fs)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout.Duration())
	defer cancel()

	return c.migrator.Ping(ctx, bytes.String(source), bytes.String(url))
}
