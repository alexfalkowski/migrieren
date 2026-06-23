package migrate

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

func (s *Migrator) sourceAndURL(ctx context.Context, db string) (context.Context, []byte, []byte, error) {
	d, err := s.config.Database(db)
	if err != nil {
		err = fmt.Errorf("%s: %w", db, err)
		return ctx, nil, nil, diagnostics.Error(ctx, err, nil)
	}

	source, err := d.GetSource(s.fs)
	if err != nil {
		ctx = diagnostics.WithStage(ctx, diagnostics.StageSource)
		return ctx, nil, nil, diagnostics.Error(ctx, err, nil)
	}

	url, err := d.GetURL(s.fs)
	if err != nil {
		ctx = diagnostics.WithStage(ctx, diagnostics.StageURL)
		return ctx, nil, nil, diagnostics.Error(ctx, err, nil)
	}

	return ctx, source, url, nil
}
