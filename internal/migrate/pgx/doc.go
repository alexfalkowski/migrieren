// Package pgx adapts golang-migrate's pgx v5 database driver.
//
// It parses Migrieren's supported pgx5 URL query parameters into the upstream
// driver config and exposes the narrow driver-construction surface used by the
// database package.
package pgx
