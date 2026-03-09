// Package database opens and instruments migrate database drivers.
//
// The package converts service-level database URLs into configured
// github.com/golang-migrate/migrate/v4/database.Driver instances.
//
// Supported schemes:
//   - pgx5://
//
// For pgx5 URLs, the package:
//   - rewrites the scheme to postgres:// for the underlying SQL driver,
//   - creates an otelsql-instrumented database handle, and
//   - registers SQL DB stats metrics.
//
// Errors are normalized to [ErrInvalidURL] and [ErrUnsupportedDriver] for
// caller stability.
//
// Telemetry setup uses runtime.Must and therefore panics if otelsql open or
// metric registration fails.
package database
