// Package health wires service health checks and health endpoint observers.
//
// The package registers:
//   - baseline noop and online checks,
//   - one migration checker per configured database, and
//   - HTTP/gRPC health observers for service endpoints.
//
// Migration health checks delegate to internal/migrate and validate that each
// configured source/database pair can be opened and inspected.
//
// Duration and timeout values are parsed from [Config] using
// go-service/time.MustParseDuration and therefore panic on invalid duration
// strings.
package health
