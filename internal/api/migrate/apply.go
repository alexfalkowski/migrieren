package migrate

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

// ApplyMigrations applies all pending up migrations for the named database.
//
// The database is resolved from configuration by name. Its migration source and
// database URL are read via the filesystem abstraction, then passed to the core
// migrator.
//
// Returns the input context, or a derived context when migration setup or the
// core migrator adds metadata, plus the resulting migration version and
// migration logs. If the database name does not exist in the configuration,
// this returns an error that wraps
// `internal/migrate.ErrNotFound` (detectable via [IsNotFound]).
//
// Source and URL resolution failures return the underlying filesystem error and
// a derived context containing a diagnostic stage. Staged failures are
// classified as "invalid_config" for transport diagnostics. Failed migrations
// may carry safe failure diagnostics on the returned error.
func (s *Migrator) ApplyMigrations(ctx context.Context, db string) (context.Context, uint64, []string, error) {
	ctx, source, url, err := s.sourceAndURL(ctx, db)
	if err != nil {
		return ctx, 0, nil, err
	}

	ctx, version, logs, err := s.migrator.ApplyMigrations(ctx, bytes.String(source), bytes.String(url))
	if err != nil {
		err = diagnostics.Error(ctx, err, logs)
	}

	return ctx, version, logs, err
}
