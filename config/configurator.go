package config

import (
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/migrieren/client"
	"github.com/alexfalkowski/migrieren/health"
	"github.com/alexfalkowski/migrieren/migrate"
)

// NewConfigurator for config.
func NewConfigurator() config.Configurator {
	cfg := &Config{}

	return cfg
}

func clientConfig(cfg config.Configurator) *client.Config {
	return &cfg.(*Config).Client
}

func healthConfig(cfg config.Configurator) *health.Config {
	return &cfg.(*Config).Health
}

func migrateConfig(cfg config.Configurator) *migrate.Config {
	return &cfg.(*Config).Migrate
}
