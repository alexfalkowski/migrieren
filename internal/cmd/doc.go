// Package cmd wires CLI commands for the service binary.
//
// The package exposes:
//   - [Module], the dependency graph used by the server command.
//   - [RegisterServer], which registers the "server" CLI command on the
//     provided commander.
//
// The server command boots the full service stack (configuration, health,
// transports, and runtime server infrastructure).
package cmd
