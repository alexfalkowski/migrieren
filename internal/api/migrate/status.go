package migrate

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// Status reports the current migration version state for the named database.
//
// The database is resolved from configuration by name. Only its database URL is
// read via the filesystem abstraction; source resolution is not needed for this
// status inspection path.
//
// If the database name does not exist in the configuration, this returns an
// error that wraps `internal/migrate.ErrNotFound` (detectable via [IsNotFound]).
// URL resolution failures return the underlying filesystem error and a derived
// context containing a URL diagnostic stage.
func (s *Migrator) Status(ctx context.Context, db string) (context.Context, *migrate.Status, error) {
	d, err := s.config.Database(db)
	if err != nil {
		return ctx, nil, fmt.Errorf("%s: %w", db, err)
	}

	url, err := d.GetURL(s.fs)
	if err != nil {
		ctx = diagnostics.WithStage(ctx, diagnostics.StageURL)
		return ctx, nil, err
	}

	return s.migrator.Status(ctx, bytes.String(url))
}
