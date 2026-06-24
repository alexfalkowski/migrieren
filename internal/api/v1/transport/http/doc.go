// Package http registers the HTTP RPC facade for the migrieren.v1 API.
//
// The shared HTTP RPC runtime exposes the generated gRPC service methods at the
// repository-documented HTTP routes while delegating migration behavior to the
// versioned API migrator.
package http
