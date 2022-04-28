package config

import (
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/migrieren/health"
)

// NewConfigurator for config.
func NewConfigurator() config.Configurator {
	cfg := &Config{}

	return cfg
}

func healthConfig(cfg config.Configurator) *health.Config {
	return &cfg.(*Config).Health
}
