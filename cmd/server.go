package cmd

import (
	ac "github.com/alexfalkowski/auth/client"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/migrieren/config"
	"github.com/alexfalkowski/migrieren/server/health"
	v1 "github.com/alexfalkowski/migrieren/server/v1"
	"github.com/alexfalkowski/migrieren/transport"
	"go.uber.org/fx"
)

// ServerOptions for cmd.
var ServerOptions = []fx.Option{
	fx.NopLogger, runtime.Module, debug.Module, feature.Module,
	telemetry.Module, health.Module, metrics.Module,
	Module, config.Module, transport.Module,
	v1.Module, ac.Module,
}
