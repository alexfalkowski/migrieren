package cmd

import (
	ac "github.com/alexfalkowski/auth/client"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	mc "github.com/alexfalkowski/migrieren/client"
	"github.com/alexfalkowski/migrieren/config"
	"github.com/alexfalkowski/migrieren/transport"
	"go.uber.org/fx"
)

// ClientOptions for cmd.
var ClientOptions = []fx.Option{
	fx.NopLogger, runtime.Module, debug.Module, feature.Module,
	telemetry.Module, metrics.Module,
	Module, config.Module, transport.Module,
	ac.Module, mc.CommandModule,
}
