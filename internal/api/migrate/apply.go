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
func (s *Migrator) ApplyMigrations(ctx context.Context, db string) (context.Context, uint64, []string, error) {
	source, url, err := s.sourceAndURL(db)
	if err != nil {
		return ctx, 0, nil, err
	}

	ctx, version, logs, err := s.migrator.ApplyMigrations(ctx, bytes.String(source), bytes.String(url))
	if err != nil {
		return ctx, version, logs, diagnostics.Error(err, logs)
	}

	return ctx, version, logs, nil
}
