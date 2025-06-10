package config

import (
	"github.com/alexfalkowski/go-service/v2/config"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module for fx.
var Module = di.Module(
	di.Constructor(config.NewConfig[Config]),
	di.Decorate(decorateConfig),
	di.Constructor(healthConfig),
	di.Constructor(migrateConfig),
)
