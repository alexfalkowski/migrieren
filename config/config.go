package config

import (
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/migrieren/health"
	"github.com/alexfalkowski/migrieren/migrate"
)

// Config for the service.
type Config struct {
	Health         *health.Config  `yaml:"health,omitempty" json:"health,omitempty" toml:"health,omitempty"`
	Migrate        *migrate.Config `yaml:"migrate,omitempty" json:"migrate,omitempty" toml:"migrate,omitempty"`
	*config.Config `yaml:",inline" json:",inline" toml:",inline"`
}

func decorateConfig(cfg *Config) *config.Config {
	return cfg.Config
}

func healthConfig(cfg *Config) *health.Config {
	return cfg.Health
}

func migrateConfig(cfg *Config) *migrate.Config {
	return cfg.Migrate
}
