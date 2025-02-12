package health

import (
	"context"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/health"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/migrieren/internal/migrate"
	"github.com/alexfalkowski/migrieren/internal/migrate/migrator"
	"go.uber.org/fx"
)

// Params for health.
type Params struct {
	fx.In

	Health   *Config
	Migrate  *migrate.Config
	Migrator migrator.Migrator
}

// NewRegistrations for health.
func NewRegistrations(params Params) (health.Registrations, error) {
	d := time.MustParseDuration(params.Health.Duration)
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
	db       *migrate.Database
	migrator migrator.Migrator
}

// Check the migrator.
func (c *migratorChecker) Check(ctx context.Context) error {
	u, err := c.db.GetURL()
	if err != nil {
		return err
	}

	return c.migrator.Ping(ctx, c.db.Source, u)
}
