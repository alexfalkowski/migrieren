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
		config.ConfigModule,
		marshaller.Module,
		fx.Provide(v1ClientConfig),
		fx.Provide(healthConfig),
		fx.Provide(migrateConfig),
	)
)
