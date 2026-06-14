// Package http registers the HTTP RPC facade for the migrieren.v1 API.
//
// The shared HTTP RPC runtime exposes the generated gRPC Migrate method at the
// repository-documented HTTP route while keeping request handling delegated to
// the gRPC transport implementation.
package http
