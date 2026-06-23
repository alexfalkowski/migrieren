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

	// ErrMigrationCanceled is returned when migration execution stops because
	// the request context was canceled.
	ErrMigrationCanceled = errors.New("migration canceled")

	// ErrMigrationDeadlineExceeded is returned when migration execution stops
	// because the request context deadline expired.
	ErrMigrationDeadlineExceeded = errors.New("migration deadline exceeded")

	// ErrInvalidPing happened.
	ErrInvalidPing = errors.New("invalid ping")
)

// NewMigrator creates a new [Migrator] instance.
//
// The returned migrator is stateless; each call to [Migrator.Migrate],
// [Migrator.ApplyMigrations], [Migrator.Ping], or [Migrator.Status] creates
// fresh underlying resources.
// Migration resources are closed before returning on normal completion and
// closed asynchronously when an in-flight migration is stopped by context
// cancellation or deadline expiry.
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
//   - version: the target migration version. Public API callers must reject
//     zero before calling Migrate; this method forwards the value to
//     golang-migrate.
//
// Output:
//   - ctx: the input context, or a derived context containing "migrateError" when
//     source/database setup or migration execution fails.
//   - logs: in-memory migration logs captured during the operation. Logs may be
//     returned on successful migrations, no-op migrations, and migration
//     execution failures. At most 100 entries are returned; truncation is marked
//     by "migration logs truncated" as the first entry.
//   - error: nil on success; otherwise one of:
//     [ErrInvalidConfig] (cannot open src/db),
//     [ErrInvalidMigration] (migration failed),
//     [ErrMigrationCanceled] (request context canceled),
//     [ErrMigrationDeadlineExceeded] (request context deadline expired).
//
// The underlying migrate.ErrNoChange is treated as a successful no-op; logs are
// still returned.
func (m *Migrator) Migrate(ctx context.Context, src, db string, version uint64) (context.Context, []string, error) {
	if err := ctx.Err(); err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("migrateError", meta.Error(err)))
		return ctx, nil, migrationError(err)
	}

	migrator, err := m.newMigrator(ctx, src, db)
	if err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("migrateError", meta.Error(err)))
		return ctx, nil, ErrInvalidConfig
	}

	logger := logger.New()
	migrator.Log = logger

	if err := m.migrate(ctx, migrator, version); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return ctx, logger.Logs(), nil
		}

		ctx = meta.WithAttributes(ctx, meta.NewPair("migrateError", meta.Error(err)))
		return ctx, logger.Logs(), migrationError(err)
	}

	return ctx, logger.Logs(), nil
}

// ApplyMigrations applies all pending up migrations for the database identified
// by db.
//
// Inputs:
//   - ctx: a service context used for metadata/telemetry.
//   - src: migration source URL (for example "file://...").
//   - db: database URL (for example a Postgres URL).
//
// Output:
//   - ctx: the input context, or a derived context containing "migrateError" when
//     source/database setup, migration execution, or final version inspection
//     fails.
//   - version: the resulting clean migration version. It is zero when no
//     migration version has been recorded.
//   - logs: in-memory migration logs captured during the operation. Logs may be
//     returned on successful migrations, no-op migrations, and migration
//     execution failures. At most 100 entries are returned; truncation is marked
//     by "migration logs truncated" as the first entry.
//   - error: nil on success; otherwise one of:
//     [ErrInvalidConfig] (cannot open src/db),
//     [ErrInvalidMigration] (migration failed or final version cannot be read),
//     [ErrMigrationCanceled] (request context canceled),
//     [ErrMigrationDeadlineExceeded] (request context deadline expired).
//
// The underlying migrate.ErrNoChange is treated as a successful no-op; logs and
// the current migration version are still returned. As with [Migrator.Status],
// strict request cancellation depends on upstream migrate v4 context support in
// database driver paths.
func (m *Migrator) ApplyMigrations(ctx context.Context, src, db string) (context.Context, uint64, []string, error) {
	if err := ctx.Err(); err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("migrateError", meta.Error(err)))
		return ctx, 0, nil, migrationError(err)
	}

	migrator, err := m.newMigrator(ctx, src, db)
	if err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("migrateError", meta.Error(err)))
		return ctx, 0, nil, ErrInvalidConfig
	}

	logger := logger.New()
	migrator.Log = logger

	version, err := m.applyMigrations(ctx, migrator)
	if err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("migrateError", meta.Error(err)))
		return ctx, 0, logger.Logs(), migrationError(err)
	}

	return ctx, version, logger.Logs(), nil
}

// Ping validates that the database can be opened and reached.
//
// Ping does not apply any migrations. It pings the database with ctx instead
// of using golang-migrate's Version path, because the pgx migrate driver does
// not accept a request context for Version. Ping does not inspect or open a
// migration source.
//
// Returns the input context, or a derived context containing "pingError" when
// database configuration or connectivity inspection fails, plus nil on success
// or one of:
//   - [ErrInvalidConfig] if db cannot be opened.
//   - [ErrInvalidPing] if the database cannot be pinged.
func (m *Migrator) Ping(ctx context.Context, db string) (context.Context, error) {
	if err := database.Ping(ctx, db); err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("pingError", meta.Error(err)))
		if errors.Is(err, database.ErrInvalidURL) || errors.Is(err, database.ErrUnsupportedDriver) {
			return ctx, ErrInvalidConfig
		}

		return ctx, ErrInvalidPing
	}

	return ctx, nil
}

func (m *Migrator) newMigrator(ctx context.Context, src, db string) (*migrate.Migrate, error) {
	s, err := source.Open(src)
	if err != nil {
		return nil, err
	}

	d, err := database.Open(ctx, db)
	if err != nil {
		_ = s.Close()
		return nil, err
	}

	migrator, err := migrate.NewWithInstance(src, s, db, d)
	if err != nil {
		_ = s.Close()
		_ = d.Close()

		return nil, err
	}

	return migrator, nil
}

func (m *Migrator) migrate(ctx context.Context, migrator *migrate.Migrate, version uint64) error {
	return m.run(ctx, migrator, func() error {
		return migrator.Migrate(uint(version))
	})
}

func (m *Migrator) applyMigrations(ctx context.Context, migrator *migrate.Migrate) (uint64, error) {
	var version uint64
	err := m.run(ctx, migrator, func() error {
		if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}

		current, _, err := migrator.Version()
		if errors.Is(err, migrate.ErrNilVersion) {
			return nil
		}
		if err != nil {
			return err
		}

		version = uint64(current)

		return nil
	})

	return version, err
}

func (m *Migrator) run(ctx context.Context, migrator *migrate.Migrate, operation func() error) error {
	if err := ctx.Err(); err != nil {
		_, _ = migrator.Close()
		return err
	}

	result := make(chan error, 1)

	go func() {
		result <- operation()
	}()

	select {
	case err := <-result:
		_, _ = migrator.Close()
		return err
	case <-ctx.Done():
		// Close can wait behind an in-flight statement, so cancellation must
		// return promptly and leave driver cleanup to finish asynchronously.
		select {
		case migrator.GracefulStop <- true:
		default:
		}
		go func() {
			_, _ = migrator.Close()
		}()

		return ctx.Err()
	}
}

func migrationError(err error) error {
	switch {
	case errors.Is(err, context.Canceled):
		return ErrMigrationCanceled
	case errors.Is(err, context.DeadlineExceeded):
		return ErrMigrationDeadlineExceeded
	default:
		return ErrInvalidMigration
	}
}
