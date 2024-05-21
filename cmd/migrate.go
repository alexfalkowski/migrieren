package cmd

import (
	"github.com/alexfalkowski/go-service/compressor"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	mc "github.com/alexfalkowski/migrieren/client"
	"github.com/alexfalkowski/migrieren/cmd/migrate"
	"github.com/alexfalkowski/migrieren/config"
	"go.uber.org/fx"
)

// MigrateOptions for cmd.
var MigrateOptions = []fx.Option{
	compressor.Module, marshaller.Module,
	telemetry.Module, metrics.Module,
	config.Module, mc.Module, migrate.Module, Module,
}
