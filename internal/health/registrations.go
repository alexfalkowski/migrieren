package health

import (
	"context"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/health"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/migrieren/internal/migrate"
	"github.com/alexfalkowski/migrieren/internal/migrate/migrator"
	"go.uber.org/fx"
)

// Params for health.
type Params struct {
	fx.In

	Health   *Config
	Migrate  *migrate.Config
	FS       *os.FS
	Migrator migrator.Migrator
}

// NewRegistrations for health.
func NewRegistrations(params Params) (health.Registrations, error) {
	d := time.MustParseDuration(params.Health.Duration)
	registrations := health.Registrations{
		server.NewRegistration("noop", d, checker.NewNoopChecker()),
		server.NewOnlineRegistration(d, d),
	}

	for _, db := range params.Migrate.Databases {
		checker := &migratorChecker{db: db, fs: params.FS, migrator: params.Migrator}
		reg := server.NewRegistration(db.Name, d, checker)

		registrations = append(registrations, reg)
	}

	return registrations, nil
}

type migratorChecker struct {
	db       *migrate.Database
	fs       *os.FS
	migrator migrator.Migrator
}

// Check the migrator.
func (c *migratorChecker) Check(ctx context.Context) error {
	source, err := c.db.GetSource(c.fs)
	if err != nil {
		return err
	}

	url, err := c.db.GetURL(c.fs)
	if err != nil {
		return err
	}

	return c.migrator.Ping(ctx, bytes.String(source), bytes.String(url))
}
