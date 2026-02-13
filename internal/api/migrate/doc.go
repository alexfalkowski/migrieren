// Package migrate provides a transport-facing adapter around the core migration
// engine.
//
// This package is used by API/transport layers (for example gRPC or HTTP) to
// execute migrations by logical database name, based on the service
// configuration.
//
// # Responsibilities
//
// The adapter is responsible for:
//
//   - Looking up a database entry from the configured [github.com/alexfalkowski/migrieren/internal/migrate.Config].
//   - Resolving the migration source and database URL for that entry via the
//     service filesystem abstraction (see github.com/alexfalkowski/go-service/v2/os.FS).
//   - Delegating the actual migration execution to the core migrator
//     (see [github.com/alexfalkowski/migrieren/internal/migrate.Migrator]).
//
// The transport-facing API intentionally does not expose the underlying
// source/URL resolution details to callers; callers typically provide only a
// database name and a target version.
//
// # Not found handling
//
// When a database name is not present in the configuration,
// [Migrator.Migrate] returns a wrapped error that includes the original sentinel
// error [github.com/alexfalkowski/migrieren/internal/migrate.ErrNotFound]. Use
// [IsNotFound] to test for this condition when mapping errors to transport
// status codes.
//
// # Observability
//
// Underlying migration and driver errors are handled by the core migrator and
// are typically attached to request metadata for telemetry. This adapter focuses
// on configuration lookup and secret/source resolution.
package migrate
