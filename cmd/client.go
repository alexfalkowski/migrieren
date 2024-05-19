package cmd

import (
	"github.com/alexfalkowski/go-service/compressor"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	mc "github.com/alexfalkowski/migrieren/client"
	"github.com/alexfalkowski/migrieren/config"
	"go.uber.org/fx"
)

// ClientOptions for cmd.
var ClientOptions = []fx.Option{
	runtime.Module, feature.Module,
	compressor.Module, marshaller.Module,
	telemetry.Module, metrics.Module,
	config.Module, mc.CommandModule, Module,
}
