package migrate

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// Plan reports current status and pending up migration versions for the named
// database without applying migrations.
//
// The database is resolved from configuration by name. Its migration source and
// database URL are read via the filesystem abstraction, then passed to the core
// migrator for non-mutating inspection.
func (s *Migrator) Plan(ctx context.Context, db string) (context.Context, *migrate.Plan, error) {
	source, url, err := s.sourceAndURL(db)
	if err != nil {
		return ctx, nil, err
	}

	ctx, plan, err := s.migrator.Plan(ctx, bytes.String(source), bytes.String(url))
	if err != nil {
		return ctx, plan, diagnostics.Error(err, nil)
	}

	return ctx, plan, nil
}
