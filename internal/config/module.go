package config

import (
	"github.com/alexfalkowski/go-service/v2/config"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module wires configuration loading and extraction helpers into the DI
// container.
var Module = di.Module(
	di.Constructor(config.NewConfig[Config]),
	di.Decorate(decorateConfig),
	di.Constructor(healthConfig),
	di.Constructor(migrateConfig),
)
