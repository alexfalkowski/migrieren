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

// NewMigrator checker.
func NewMigrator(db *migrate.Database, fs *os.FS, migrator *migrate.Migrator, timeout time.Duration) *Migrator {
	return &Migrator{db: db, fs: fs, migrator: migrator, timeout: timeout}
}

// Migrator checker.
type Migrator struct {
	db       *migrate.Database
	fs       *os.FS
	migrator *migrate.Migrator
	timeout  time.Duration
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

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.migrator.Ping(ctx, bytes.String(source), bytes.String(url))
}
