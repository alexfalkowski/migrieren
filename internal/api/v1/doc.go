// Package v1 wires the versioned API surface into the service dependency graph.
//
// The package owns DI composition for the v1 transport adapters. The protobuf
// contract itself lives in api/migrieren/v1, while transport-specific behavior
// lives under the transport subpackages.
package v1
