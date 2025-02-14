package cmd

import (
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/module"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/transport"
	v1 "github.com/alexfalkowski/migrieren/internal/api/v1"
	"github.com/alexfalkowski/migrieren/internal/config"
	"github.com/alexfalkowski/migrieren/internal/health"
)

// RegisterServer for cmd.
func RegisterServer(command *cmd.Command) {
	command.AddServer("server", "Start migrieren server",
		module.Module, debug.Module, feature.Module,
		telemetry.Module, transport.Module, health.Module,
		config.Module, v1.Module, cmd.Module,
	)
}
