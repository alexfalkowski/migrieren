package migrate

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/migrieren/internal/migrate/database"
	"github.com/alexfalkowski/migrieren/internal/migrate/source"
	"github.com/alexfalkowski/migrieren/internal/migrate/telemetry/logger"
	"github.com/golang-migrate/migrate/v4"
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
func NewMigrator() *Migrator {
	return &Migrator{}
}

// Migrator using migrate.
type Migrator struct{}

// Migrate a database to a version and returning the database logs.
func (m *Migrator) Migrate(ctx context.Context, src, db string, version uint64) ([]string, error) {
	s, err := source.Open(src)
	if err != nil {
		meta.WithAttribute(ctx, "migrateError", meta.Error(err))

		return nil, ErrInvalidConfig
	}

	d, err := database.Open(db)
	if err != nil {
		meta.WithAttribute(ctx, "migrateError", meta.Error(err))

		return nil, ErrInvalidConfig
	}

	// Error is never returned.
	migrator, _ := migrate.NewWithInstance(src, s, db, d)

	logger := logger.New()
	migrator.Log = logger

	if err := migrator.Migrate(uint(version)); err != nil {
		meta.WithAttribute(ctx, "migrateError", meta.Error(err))

		if errors.Is(err, migrate.ErrNoChange) {
			return logger.Logs(), m.close(migrator, nil)
		}

		return logger.Logs(), m.close(migrator, ErrInvalidMigration)
	}

	return logger.Logs(), m.close(migrator, nil)
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
