package config

import (
	"github.com/alexfalkowski/go-service/config"
	v1 "github.com/alexfalkowski/migrieren/client/v1/config"
	"github.com/alexfalkowski/migrieren/health"
	"github.com/alexfalkowski/migrieren/migrate"
)

// NewConfigurator for config.
func NewConfigurator() config.Configurator {
	cfg := &Config{}

	return cfg
}

func v1ClientConfig(cfg config.Configurator) *v1.Config {
	return &cfg.(*Config).Client.V1
}

func healthConfig(cfg config.Configurator) *health.Config {
	return &cfg.(*Config).Health
}

func migrateConfig(cfg config.Configurator) *migrate.Config {
	return &cfg.(*Config).Migrate
}
