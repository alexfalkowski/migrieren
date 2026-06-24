package migrate

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

// Migrate migrates the named database to the given target version.
//
// The database is resolved from configuration by name. Its migration source and
// database URL are read via the filesystem abstraction, then passed to the core
// migrator.
func (s *Migrator) Migrate(ctx context.Context, db string, version uint64) (context.Context, []string, error) {
	source, url, err := s.sourceAndURL(db)
	if err != nil {
		return ctx, nil, err
	}

	ctx, logs, err := s.migrator.Migrate(ctx, bytes.String(source), bytes.String(url), version)
	if err != nil {
		return ctx, logs, diagnostics.Error(err, logs)
	}

	return ctx, logs, nil
}
