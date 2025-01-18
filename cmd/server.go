package cmd

import (
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/module"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/transport"
	v1 "github.com/alexfalkowski/migrieren/api/v1"
	"github.com/alexfalkowski/migrieren/config"
	"github.com/alexfalkowski/migrieren/health"
	"go.uber.org/fx"
)

// ServerOptions for cmd.
var ServerOptions = []fx.Option{
	module.Module, debug.Module, feature.Module,
	telemetry.Module, transport.Module, health.Module,
	config.Module, v1.Module, Module,
}
