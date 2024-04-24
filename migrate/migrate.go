package migrate

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/migrieren/migrate/migrator"
	"github.com/alexfalkowski/migrieren/migrate/telemetry/tracer"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // need this for migrations to work.
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/multierr"
)

var (
	// ErrInvalidConfig for source or db.
	ErrInvalidConfig = errors.New("invalid config")

	// ErrInvalidMigration happened.
	ErrInvalidMigration = errors.New("invalid migration")

	// ErrInvalidPing happened.
	ErrInvalidPing = errors.New("invalid ping")
)

// NewMigrator for databases.
func NewMigrator(t trace.Tracer) migrator.Migrator {
	var m migrator.Migrator = &Migrator{}
	m = tracer.NewMigrator(m, t)

	return m
}

// Migrator using migrate.
type Migrator struct{}

// Migrate a database to a version and returning the database logs.
func (m *Migrator) Migrate(ctx context.Context, source, db string, version uint64) ([]string, error) {
	mig, err := migrate.New(source, db)
	if err != nil {
		meta.WithAttribute(ctx, "migrateError", meta.Error(err))

		return nil, ErrInvalidConfig
	}

	logger := &logger{logs: make([]string, 0)}
	mig.Log = logger

	if err := mig.Migrate(uint(version)); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Printf(err.Error())

			return logger.logs, m.close(mig, nil)
		}

		meta.WithAttribute(ctx, "migrateError", meta.Error(err))

		return logger.logs, m.close(mig, ErrInvalidMigration)
	}

	return logger.logs, m.close(mig, nil)
}

// Ping the migrator.
func (m *Migrator) Ping(ctx context.Context, source, db string) error {
	mig, err := migrate.New(source, db)
	if err != nil {
		meta.WithAttribute(ctx, "pingError", meta.Error(err))

		return ErrInvalidConfig
	}

	if _, _, err := mig.Version(); err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		meta.WithAttribute(ctx, "pingError", meta.Error(err))

		return ErrInvalidPing
	}

	return nil
}

func (m *Migrator) close(mig *migrate.Migrate, err error) error {
	errSource, errDB := mig.Close()

	return multierr.Combine(errSource, errDB, err)
}
