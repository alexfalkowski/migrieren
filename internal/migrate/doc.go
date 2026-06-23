// Package migrate provides the core migration engine used by the service.
//
// This package wraps github.com/golang-migrate/migrate/v4 to perform database
// schema migrations and to "ping" databases to verify connectivity.
//
// # Overview
//
// The primary type is [Migrator], which offers:
//
//   - [Migrator.Migrate] to open a migration source and database, migrate the
//     database to a target version, and return in-memory migration logs.
//   - [Migrator.Ping] to validate that the database can be opened and reached.
//   - [Migrator.Status] to report the current migration version and dirty state.
//
// The migrator operates on two string inputs:
//
//   - src: a migration source URL (for example a file-based source such as
//     "file://..."; supported source drivers are registered by the internal
//     source wiring).
//   - db: a database URL (for example a Postgres URL; supported database drivers
//     are registered by the internal database wiring).
//
// # Error model
//
// This package intentionally hides underlying driver errors behind a small set
// of stable, exported sentinel errors:
//
//   - [ErrInvalidConfig] is returned when either the migration source or the
//     database URL cannot be opened.
//   - [ErrInvalidMigration] is returned when a migration attempt fails for
//     reasons other than "no change".
//   - [ErrMigrationCanceled] and [ErrMigrationDeadlineExceeded] are returned
//     when migration execution stops because the request context is canceled or
//     its deadline expires.
//   - [ErrInvalidPing] is returned when pinging/inspecting the database fails.
//   - [ErrInvalidStatus] is returned when migration status inspection fails.
//
// Underlying errors are attached to the provided context as metadata attributes
// for observability, but are not returned directly to callers.
//
// # Logs
//
// Migration logs are captured in memory and returned from [Migrator.Migrate] as
// a slice of strings. The logger keeps at most 100 entries; when output exceeds
// that cap, the first returned entry is "migration logs truncated" and the
// remaining entries are the latest log lines. When the underlying migrate
// engine reports migrate.ErrNoChange, it is treated as a successful no-op and
// the accumulated logs are still returned.
//
// # Resource management
//
// The underlying migrate engine allocates resources for both the migration
// source and the database connection. Normal completion closes those resources
// before returning to the caller. When the request context is canceled or its
// deadline expires while a migration is running, this package requests a
// graceful stop and closes resources asynchronously so cancellation can return
// promptly.
//
// Status inspection is non-migrating, but the underlying migrate v4 database
// driver version path does not accept a request context. The migrator checks
// cancellation before and after status inspection, but cannot interrupt every
// upstream driver path until migrate v5 provides context-aware driver APIs.
package migrate
