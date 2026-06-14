// Package database opens and checks migrate database drivers.
//
// The package owns supported database-driver selection for migration and health
// paths. It currently accepts pgx5 URLs, rewrites them for the underlying
// Postgres SQL driver, applies request-deadline statement timeouts, and wires
// database telemetry around opened connections.
package database
