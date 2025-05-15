package migrate

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexfalkowski/go-service/bytes"
	"github.com/alexfalkowski/migrieren/internal/migrate"
	"github.com/alexfalkowski/migrieren/internal/migrate/migrator"
)

// ErrNotFound for service.
var ErrNotFound = errors.New("not found")

// IsNotFound for service.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// NewMigrator for the different transports.
func NewMigrator(cfg *migrate.Config, mig migrator.Migrator) *Migrator {
	return &Migrator{config: cfg, migrator: mig}
}

// Migrator for the different transports.
type Migrator struct {
	migrator migrator.Migrator
	config   *migrate.Config
}

// Migrate the database.
func (s *Migrator) Migrate(ctx context.Context, db string, version uint64) ([]string, error) {
	d := s.config.Database(db)
	if d == nil {
		return nil, fmt.Errorf("%s: %w", db, ErrNotFound)
	}

	source, err := d.GetSource()
	if err != nil {
		return nil, err
	}

	url, err := d.GetURL()
	if err != nil {
		return nil, err
	}

	return s.migrator.Migrate(ctx, bytes.String(source), bytes.String(url), version)
}
