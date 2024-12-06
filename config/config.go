package config

import (
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/migrieren/health"
	"github.com/alexfalkowski/migrieren/migrate"
)

// NewConfigurator for config.
func NewConfig(i *cmd.InputConfig) (*Config, error) {
	c := &Config{}

	return c, i.Decode(c)
}

// IsEnabled for config.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for the service.
type Config struct {
	Health         *health.Config  `yaml:"health,omitempty" json:"health,omitempty" toml:"health,omitempty"`
	Migrate        *migrate.Config `yaml:"migrate,omitempty" json:"migrate,omitempty" toml:"migrate,omitempty"`
	*config.Config `yaml:",inline" json:",inline" toml:",inline"`
}

func decorateConfig(cfg *Config) *config.Config {
	if !IsEnabled(cfg) {
		return nil
	}

	return cfg.Config
}

func healthConfig(cfg *Config) *health.Config {
	if !IsEnabled(cfg) {
		return nil
	}

	return cfg.Health
}

func migrateConfig(cfg *Config) *migrate.Config {
	if !IsEnabled(cfg) {
		return nil
	}

	return cfg.Migrate
}
