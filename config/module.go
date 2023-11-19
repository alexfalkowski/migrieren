package config

import (
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/go-service/marshaller"
	"go.uber.org/fx"
)

var (
	// Module for fx.
	Module = fx.Options(
		fx.Provide(NewConfigurator),
		marshaller.Module,
		config.ConfigModule,
		fx.Provide(v1ClientConfig), fx.Provide(v1AuthClientConfig),
		fx.Provide(healthConfig), fx.Provide(migrateConfig),
	)
)
