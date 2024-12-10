package config

import (
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/migrieren/token"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	token.Module,
	fx.Provide(config.NewConfig[Config]),
	config.Module,
	fx.Decorate(decorateConfig),
	fx.Provide(healthConfig),
	fx.Provide(migrateConfig),
)
