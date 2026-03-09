// Package grpc provides the API v1 gRPC transport.
//
// It implements the protobuf service defined in api/migrieren/v1 and delegates
// migration execution to the transport-facing migrator adapter.
//
// Error mapping:
//   - internal/api/migrate not-found errors are mapped to gRPC NotFound.
//   - all other errors are mapped to gRPC Internal.
//
// Responses include migration logs and request metadata copied from the
// operation context.
package grpc
