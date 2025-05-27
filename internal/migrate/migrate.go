package migrate

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/migrieren/internal/migrate/migrator"
	"github.com/alexfalkowski/migrieren/internal/migrate/telemetry/logger"
	tt "github.com/alexfalkowski/migrieren/internal/migrate/telemetry/tracer"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // need this for migrations to work.
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
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
func NewMigrator(tracer *tracer.Tracer) migrator.Migrator {
	var m migrator.Migrator = &Migrator{}
	m = tt.NewMigrator(m, tracer)

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

	logger := logger.New()
	mig.Log = logger

	if err := mig.Migrate(uint(version)); err != nil {
		meta.WithAttribute(ctx, "migrateError", meta.Error(err))

		if errors.Is(err, migrate.ErrNoChange) {
			return logger.Logs(), m.close(mig, nil)
		}

		return logger.Logs(), m.close(mig, ErrInvalidMigration)
	}

	return logger.Logs(), m.close(mig, nil)
}

// Ping the migrator.
func (m *Migrator) Ping(ctx context.Context, source, db string) error {
	_, err := migrate.New(source, db)
	if err != nil {
		meta.WithAttribute(ctx, "pingError", meta.Error(err))

		return ErrInvalidConfig
	}

	return nil
}

func (m *Migrator) close(mig *migrate.Migrate, err error) error {
	errSource, errDB := mig.Close()

	return errors.Join(errSource, errDB, err)
}
