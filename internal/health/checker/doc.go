// Package checker provides health checker implementations used by
// internal/health.
//
// It exposes:
//   - a noop checker constructor for baseline liveness checks,
//   - a migration checker that validates database/source reachability by calling
//     the core migrator Ping API.
//
// The migration checker resolves source and URL values through go-service os.FS
// and applies a per-check timeout context.
package checker
