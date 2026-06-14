// Package checker contains health check implementations for Migrieren.
//
// It provides a noop checker wrapper and a migration-target checker that
// resolves configured source/database references, validates source syntax or
// openness, and pings the configured database within the health timeout.
package checker
