// Package redis wires the Redis-based synchronization client used by the
// migration service.
//
// The package exposes a constructor that:
//   - reads a Redis connection URL from a configurable source path,
//   - builds a github.com/redis/go-redis/v9 client,
//   - disables Redis maintenance notifications, and
//   - instruments Redis operations with OpenTelemetry tracing and metrics.
//
// The returned client type is github.com/go-redsync/redsync/v4, which is used
// by the migration layer to coordinate distributed locks across service
// instances.
//
// Configuration:
//   - [Config.URL] is resolved via go-service's os.FS source resolution
//     (for example plain URLs or file-backed sources like "file:secrets/redis").
//
// Errors:
//   - [NewClient] returns an error when reading URL data fails or when the URL
//     cannot be parsed as a Redis connection string.
//   - Telemetry instrumentation uses runtime.Must and therefore panics if
//     tracing/metrics instrumentation cannot be attached.
package redis
