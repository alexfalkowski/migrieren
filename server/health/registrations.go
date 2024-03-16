package health

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/health"
	h "github.com/alexfalkowski/migrieren/health"
	"github.com/alexfalkowski/migrieren/migrate"
	"github.com/alexfalkowski/migrieren/migrate/migrator"
	"go.uber.org/fx"
)

// Params for health.
type Params struct {
	fx.In

	Health   *h.Config
	Migrate  *migrate.Config
	Migrator migrator.Migrator
}

// NewRegistrations for health.
func NewRegistrations(params Params) (health.Registrations, error) {
	d, err := time.ParseDuration(params.Health.Duration)
	if err != nil {
		return nil, err
	}

	registrations := health.Registrations{
		server.NewRegistration("noop", d, checker.NewNoopChecker()),
	}

	for _, db := range params.Migrate.Databases {
		checker := &migratorChecker{db: db, migrator: params.Migrator}
		reg := server.NewRegistration(db.Name, d, checker)

		registrations = append(registrations, reg)
	}

	return registrations, nil
}

type migratorChecker struct {
	db       migrate.Database
	migrator migrator.Migrator
}

// Check the migrator.
func (c *migratorChecker) Check(ctx context.Context) error {
	return c.migrator.Ping(ctx, c.db.Source, c.db.URL)
}
