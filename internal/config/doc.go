// Package config defines and wires service configuration.
//
// The package provides:
//   - [Config], the root service configuration model.
//   - [Module], DI wiring that loads [Config], exposes the embedded base
//     go-service config, and provides subsystem configs for health, migration,
//     and Redis modules.
//
// Configuration values are unmarshaled from the runtime input source selected
// by the service CLI (for example `-i file:...`).
package config
