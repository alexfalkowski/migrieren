// Package migrate provides the core migration engine used by the service.
//
// This package wraps github.com/golang-migrate/migrate/v4 to perform database
// schema migrations and to "ping" databases to verify connectivity and basic
// migrator health.
//
// # Overview
//
// The primary type is [Migrator], which offers:
//
//   - [Migrator.Migrate] to migrate a database to a target version, returning
//     in-memory migration logs.
//   - [Migrator.Ping] to validate that the migration source and database can be
//     opened and that the database is reachable.
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
//   - [ErrInvalidPing] is returned when pinging/inspecting the database fails.
//
// Underlying errors are attached to the provided context as metadata attributes
// for observability, but are not returned directly to callers.
//
// # Logs
//
// Migration logs are captured in memory and returned from [Migrator.Migrate] as
// a slice of strings. When the underlying migrate engine reports
// migrate.ErrNoChange, it is treated as a successful no-op and the accumulated
// logs are still returned.
//
// # Resource management
//
// The underlying migrate engine allocates resources for both the migration
// source and the database connection. This package ensures those resources are
// closed after each operation and joins any close errors with the operation
// error where applicable.
package migrate
