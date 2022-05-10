package migrate

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/migrieren/migrate/migrator"
	"github.com/alexfalkowski/migrieren/migrate/trace/opentracing"
	"github.com/golang-migrate/migrate/v4"

	// These are here to make sure we can use migrate. Add here to extend it.
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

var (
	// ErrInvalidConfig for source or db.
	ErrInvalidConfig = errors.New("invalid config")

	// ErrInvalidMigration happened.
	ErrInvalidMigration = errors.New("invalid migration")
)

// NewMigrator for databases.
func NewMigrator(tracer opentracing.Tracer) migrator.Migrator {
	var m migrator.Migrator = &Migrator{}
	m = opentracing.NewMigrator(m, tracer)

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

	defer mig.Close()

	logger := &logger{logs: make([]string, 0)}
	mig.Log = logger

	if err := mig.Migrate(uint(version)); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return []string{err.Error()}, nil
		}

		meta.WithAttribute(ctx, "migrate.error", err.Error())

		return nil, ErrInvalidMigration
	}

	return logger.logs, nil
}

// Ping the migrator.
func (m *Migrator) Ping(ctx context.Context, source, db string) error {
	mig, err := migrate.New(source, db)
	if err != nil {
		return err
	}

	if _, _, err := mig.Version(); err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return err
	}

	return nil
}
