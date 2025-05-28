package config

import (
	"github.com/alexfalkowski/go-service/v2/config"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(config.NewConfig[Config]),
	fx.Decorate(decorateConfig),
	fx.Provide(healthConfig),
	fx.Provide(migrateConfig),
)
