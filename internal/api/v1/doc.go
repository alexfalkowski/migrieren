// Package v1 wires API v1 transport modules.
//
// The package composes:
//   - the transport-facing migration adapter (internal/api/migrate),
//   - the core migration engine (internal/migrate),
//   - gRPC server construction and registration, and
//   - HTTP RPC route registration that forwards to the gRPC handler.
//
// Consumers typically import [Module] into the application dependency graph.
package v1
