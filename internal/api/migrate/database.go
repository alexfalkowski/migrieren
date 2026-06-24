package migrate

import (
	"fmt"

	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

func (s *Migrator) sourceAndURL(db string) ([]byte, []byte, error) {
	d, err := s.config.Database(db)
	if err != nil {
		return nil, nil, diagnostics.Error(fmt.Errorf("%s: %w", db, err), nil)
	}

	source, err := d.GetSource(s.fs)
	if err != nil {
		return nil, nil, diagnostics.InvalidConfig(err, diagnostics.StageSource)
	}

	url, err := d.GetURL(s.fs)
	if err != nil {
		return nil, nil, diagnostics.InvalidConfig(err, diagnostics.StageURL)
	}

	return source, url, nil
}
