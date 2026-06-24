// Package migrate provides an API-facing adapter around the core migration
// engine.
//
// This package is used by versioned API layers to execute migrations, inspect
// migration status, or list configured logical database names based on the
// service configuration.
//
// # Responsibilities
//
// The adapter is responsible for:
//
//   - Looking up a database entry from the configured [github.com/alexfalkowski/migrieren/internal/migrate.Config].
//   - Resolving the migration source and database URL for that entry via the
//     service filesystem abstraction (see github.com/alexfalkowski/go-service/v2/os.FS).
//   - Delegating migration execution and status inspection to the core migrator
//     (see [github.com/alexfalkowski/migrieren/internal/migrate.Migrator]).
//   - Listing configured logical database names without exposing source or URL
//     values.
//
// The API-facing contract intentionally does not expose the underlying
// source/URL resolution details to callers; callers provide only a database name,
// plus a target version for migration requests.
package migrate
