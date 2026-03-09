package config

import (
	"github.com/alexfalkowski/go-service/v2/config"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module registers configuration loading and projection for DI.
//
// It:
//   - constructs the root [Config] from runtime input,
//   - decorates it as the embedded go-service [config.Config], and
//   - provides subsystem config pointers for dependent modules.
var Module = di.Module(
	di.Constructor(config.NewConfig[Config]),
	di.Decorate(decorateConfig),
	di.Constructor(healthConfig),
	di.Constructor(migrateConfig),
	di.Constructor(redisConfig),
)
