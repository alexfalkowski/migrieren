package migrate

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// Status reports the current migration version state for the named database.
//
// The database is resolved from configuration by name. Only its database URL is
// read via the filesystem abstraction; source resolution is not needed for this
// status inspection path.
func (s *Migrator) Status(ctx context.Context, db string) (context.Context, *migrate.Status, error) {
	d, err := s.config.Database(db)
	if err != nil {
		return ctx, nil, fmt.Errorf("%s: %w", db, err)
	}

	url, err := d.GetURL(s.fs)
	if err != nil {
		return ctx, nil, diagnostics.InvalidConfig(err, diagnostics.StageURL)
	}

	ctx, status, err := s.migrator.Status(ctx, bytes.String(url))
	if err != nil {
		if stage := diagnostics.Stage(err); !strings.IsEmpty(stage) {
			return ctx, nil, diagnostics.InvalidConfig(err, stage)
		}

		return ctx, nil, err
	}

	return ctx, status, nil
}
