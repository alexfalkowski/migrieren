package checker

import (
	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/migrieren/internal/migrate"
	"github.com/alexfalkowski/migrieren/internal/migrate/source"
)

// NewNoopChecker returns a checker that always reports healthy.
func NewNoopChecker() *checker.NoopChecker {
	return checker.NewNoopChecker()
}

// NewMigrator constructs a health checker for one configured migration target.
//
// The checker resolves the migration source and database URL through fs, then
// validates the source reference and database connectivity within the provided
// timeout.
func NewMigrator(db *migrate.Database, fs *os.FS, migrator *migrate.Migrator, timeout time.Duration) *Migrator {
	return &Migrator{db: db, fs: fs, migrator: migrator, timeout: timeout}
}

// Migrator checks whether a configured migration target can resolve its source
// and ping its database successfully.
type Migrator struct {
	db       *migrate.Database
	fs       *os.FS
	migrator *migrate.Migrator
	timeout  time.Duration
}

// Check resolves the migration source and database URL, then validates the
// source reference and pings the target database within the configured timeout.
// GitHub source checks are syntax-only; the remote source is opened during
// migration execution.
func (c *Migrator) Check(ctx context.Context) error {
	src, err := c.db.GetSource(c.fs)
	if err != nil {
		return err
	}

	url, err := c.db.GetURL(c.fs)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if err := source.Check(ctx, bytes.String(src)); err != nil {
		return err
	}

	_, err = c.migrator.Ping(ctx, bytes.String(url))
	return err
}
