package cmd

import (
	"github.com/alexfalkowski/go-service/logger"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/alexfalkowski/go-service/otel"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/migrieren/client"
	"github.com/alexfalkowski/migrieren/config"
	"go.uber.org/fx"
)

// ClientOptions for cmd.
var ClientOptions = []fx.Option{
	fx.NopLogger, runtime.Module, marshaller.Module, otel.Module,
	Module, config.Module, logger.ZapModule,
	transport.GRPCModule, client.CommandModule,
}
