// Package health registers service health checks and observers.
//
// The package installs service-level online/noop checks, per-database migration
// target checks, and the intentionally narrower gRPC service health check used
// by clients of migrieren.v1.Service.
package health
