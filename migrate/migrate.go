package migrate

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/migrieren/migrate/migrator"
	"github.com/alexfalkowski/migrieren/migrate/telemetry/tracer"
	"github.com/golang-migrate/migrate/v4"

	// These are here to make sure we can use migrate. Add here to extend it.
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
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
func NewMigrator(t tracer.Tracer) migrator.Migrator {
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
		meta.WithAttribute(ctx, "migrate.error", err.Error())

		return nil, ErrInvalidConfig
	}

	logger := &logger{logs: make([]string, 0)}
	mig.Log = logger

	if err := mig.Migrate(uint(version)); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Printf(err.Error())

			return m.close(ctx, mig, logger, nil)
		}

		meta.WithAttribute(ctx, "migrate.error", err.Error())

		return m.close(ctx, mig, logger, ErrInvalidMigration)
	}

	return m.close(ctx, mig, logger, nil)
}

// Ping the migrator.
func (m *Migrator) Ping(ctx context.Context, source, db string) error {
	mig, err := migrate.New(source, db)
	if err != nil {
		meta.WithAttribute(ctx, "ping.error", err.Error())

		return ErrInvalidConfig
	}

	if _, _, err := mig.Version(); err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		meta.WithAttribute(ctx, "ping.error", err.Error())

		return ErrInvalidPing
	}

	return nil
}

func (m *Migrator) close(ctx context.Context, mig *migrate.Migrate, log *logger, err error) ([]string, error) {
	sourceErr, dbErr := mig.Close()
	if sourceErr != nil {
		meta.WithAttribute(ctx, "migrate.source_error", sourceErr.Error())

		return log.logs, sourceErr
	}

	if dbErr != nil {
		meta.WithAttribute(ctx, "migrate.db_error", dbErr.Error())

		return log.logs, dbErr
	}

	return log.logs, err
}
