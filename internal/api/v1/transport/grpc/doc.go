// Package grpc implements the migrieren.v1 gRPC transport.
//
// It registers the generated service server, delegates migration and status
// work to the transport-facing migrator, maps domain errors to gRPC status
// codes, and adds safe failure diagnostics as trailers for post-validation
// migration failures.
package grpc
