package cmd

import (
	"github.com/alexfalkowski/go-service/logger"
	"github.com/alexfalkowski/go-service/metrics"
	"github.com/alexfalkowski/go-service/trace"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/migrieren/config"
	"github.com/alexfalkowski/migrieren/server/health"
	v1 "github.com/alexfalkowski/migrieren/server/v1"
	"go.uber.org/fx"
)

// ServerOptions for cmd.
var ServerOptions = []fx.Option{
	fx.NopLogger, fx.Provide(NewVersion), config.Module, health.Module,
	logger.ZapModule, metrics.PrometheusModule,
	transport.GRPCServerModule, transport.GRPCJaegerModule,
	transport.HTTPServerModule, transport.HTTPJaegerModule,
	trace.JaegerOpenTracingModule,
	v1.Module,
}
