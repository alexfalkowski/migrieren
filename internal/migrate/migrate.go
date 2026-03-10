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

// NewMigrator creates a new [Migrator] instance.
//
// The returned migrator is stateless; each call to [Migrator.Migrate] or
// [Migrator.Ping] creates a new underlying github.com/golang-migrate/migrate/v4
// engine instance and ensures resources are closed before returning.
func NewMigrator() *Migrator {
	return &Migrator{}
}

// Migrator performs database schema migrations using golang-migrate.
//
// It provides a thin, service-oriented API that:
//   - Opens a migration source (src) and database driver (db) for each operation.
//   - Captures migration logs in memory and returns them to the caller.
//   - Maps underlying driver/migration errors onto stable sentinel errors,
//     attaching the original error to context metadata for observability.
type Migrator struct{}

// Migrate migrates the database identified by db to the given target version.
//
// Inputs:
//   - ctx: a service context used for metadata/telemetry.
//   - src: migration source URL (for example "file://...").
//   - db: database URL (for example a Postgres URL).
//   - version: the target migration version.
//
// Output:
//   - logs: in-memory migration logs captured during the operation.
//   - error: nil on success; otherwise one of:
//     [ErrInvalidConfig] (cannot open src/db),
//     [ErrInvalidMigration] (migration failed).
//
// The underlying migrate.ErrNoChange is treated as a successful no-op; logs are
// still returned.
func (m *Migrator) Migrate(ctx context.Context, src, db string, version uint64) ([]string, error) {
	migrator, err := m.newMigrator(src, db)
	if err != nil {
		meta.WithAttribute(ctx, "migrateError", meta.Error(err))
		return nil, ErrInvalidConfig
	}
	defer migrator.Close()

	logger := logger.New()
	migrator.Log = logger

	if err := migrator.Migrate(uint(version)); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return logger.Logs(), nil
		}

		meta.WithAttribute(ctx, "migrateError", meta.Error(err))
		return logger.Logs(), ErrInvalidMigration
	}

	return logger.Logs(), nil
}

// Ping validates that the migration source and database can be opened and that
// the database is reachable.
//
// Ping does not apply any migrations. Internally it opens the migrator and
// inspects the current version. A nil (unapplied) version is treated as healthy.
//
// Returns nil on success; otherwise one of:
//   - [ErrInvalidConfig] if src/db cannot be opened.
//   - [ErrInvalidPing] if the database cannot be inspected/pinged.
func (m *Migrator) Ping(ctx context.Context, src, db string) error {
	migrator, err := m.newMigrator(src, db)
	if err != nil {
		meta.WithAttribute(ctx, "pingError", meta.Error(err))
		return ErrInvalidConfig
	}
	defer migrator.Close()

	if _, _, err := migrator.Version(); err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return nil
		}

		meta.WithAttribute(ctx, "pingError", meta.Error(err))
		return ErrInvalidPing
	}

	return nil
}

func (m *Migrator) newMigrator(src, db string) (*migrate.Migrate, error) {
	s, err := source.Open(src)
	if err != nil {
		return nil, err
	}

	d, err := database.Open(db)
	if err != nil {
		return nil, err
	}

	return migrate.NewWithInstance(src, s, db, d)
}
