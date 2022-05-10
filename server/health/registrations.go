package health

import (
	"context"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/health"
	mhealth "github.com/alexfalkowski/migrieren/health"
	"github.com/alexfalkowski/migrieren/migrate"
	"github.com/alexfalkowski/migrieren/migrate/migrator"
	"go.uber.org/fx"
)

// Params for health.
type Params struct {
	fx.In

	Health   *mhealth.Config
	Migrate  *migrate.Config
	Migrator migrator.Migrator
}

// NewRegistrations for health.
func NewRegistrations(params Params) health.Registrations {
	registrations := health.Registrations{
		server.NewRegistration("noop", params.Health.Duration, checker.NewNoopChecker()),
	}

	for _, d := range params.Migrate.Databases {
		checker := &migratorChecker{config: d, migrator: params.Migrator}
		reg := server.NewRegistration(d.Name, params.Health.Duration, checker)

		registrations = append(registrations, reg)
	}

	return registrations
}

type migratorChecker struct {
	config   migrate.Database
	migrator migrator.Migrator
}

// Check the migrator.
func (c *migratorChecker) Check(ctx context.Context) error {
	return c.migrator.Ping(ctx, c.config.Source, c.config.URL)
}
