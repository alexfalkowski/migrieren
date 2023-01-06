package cmd

import (
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/logger"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/alexfalkowski/go-service/transport"
	v1 "github.com/alexfalkowski/migrieren/client/v1"
	"github.com/alexfalkowski/migrieren/config"
	"go.uber.org/fx"
)

// ClientOptions for cmd.
var ClientOptions = []fx.Option{
	fx.NopLogger,
	marshaller.Module, cmd.Module,
	fx.Provide(NewVersion), config.Module, logger.ZapModule,
	transport.GRPCModule, v1.Module,
}
