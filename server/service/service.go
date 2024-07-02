package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexfalkowski/migrieren/migrate"
	"github.com/alexfalkowski/migrieren/migrate/migrator"
)

// ErrNotFound for service.
var ErrNotFound = errors.New("not found")

// IsNotFound for service.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// NewService for the different transports.
func NewService(cfg *migrate.Config, mig migrator.Migrator) *Service {
	return &Service{config: cfg, migrator: mig}
}

// Service for the different transports.
type Service struct {
	migrator migrator.Migrator
	config   *migrate.Config
}

// Migrate the database.
func (s *Service) Migrate(ctx context.Context, db string, version uint64) ([]string, error) {
	d := s.config.Database(db)
	if d == nil {
		return nil, fmt.Errorf("%s: %w", db, ErrNotFound)
	}

	u, err := d.GetURL()
	if err != nil {
		return nil, err
	}

	return s.migrator.Migrate(ctx, d.Source, u, version)
}
