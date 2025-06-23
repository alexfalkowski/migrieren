package health

import (
	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/health"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/migrieren/internal/migrate"
	"github.com/alexfalkowski/migrieren/internal/migrate/migrator"
)

func registrations(mig migrator.Migrator, fs *os.FS, migCfg *migrate.Config, cfg *Config) (health.Registrations, error) {
	d := time.MustParseDuration(cfg.Duration)
	registrations := health.Registrations{
		server.NewRegistration("noop", d, checker.NewNoopChecker()),
		server.NewOnlineRegistration(d, d),
	}

	for _, db := range migCfg.Databases {
		checker := &migratorChecker{db: db, fs: fs, migrator: mig}
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
