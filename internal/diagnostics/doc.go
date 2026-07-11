// Package diagnostics wraps migration errors with safe diagnostic values for
// transport metadata.
//
// The package owns the repository's migration diagnostic key contract:
// migration-error, migration-log-count, migration-stage, and
// migration-log-last. These values are safe for HTTP response headers and gRPC
// trailers because they expose error categories, source/database setup or
// reference-resolution stages, and bounded log metadata rather than raw
// database URLs, source URLs, or secret values.
//
// Wrapped errors preserve their original cause for errors.Is and errors.As.
// Transport packages decide how to expose the values, while migration packages
// decide which operation failures should carry them.
package diagnostics
