package cmd

import (
	"github.com/alexfalkowski/go-service/logger"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/alexfalkowski/go-service/metrics"
	"github.com/alexfalkowski/go-service/otel"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/migrieren/config"
	"github.com/alexfalkowski/migrieren/server/health"
	v1 "github.com/alexfalkowski/migrieren/server/v1"
	"go.uber.org/fx"
)

// ServerOptions for cmd.
var ServerOptions = []fx.Option{
	fx.NopLogger, marshaller.Module, otel.Module, Module,
	config.Module, health.Module, logger.ZapModule,
	metrics.PrometheusModule, transport.Module, v1.Module,
}
