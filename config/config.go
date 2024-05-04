package config

import (
	av1 "github.com/alexfalkowski/auth/client/v1/config"
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/migrieren/auth"
	"github.com/alexfalkowski/migrieren/client"
	v1 "github.com/alexfalkowski/migrieren/client/v1/config"
	"github.com/alexfalkowski/migrieren/health"
	"github.com/alexfalkowski/migrieren/migrate"
)

// NewConfigurator for config.
func NewConfig(i *cmd.InputConfig) (*Config, error) {
	c := &Config{}

	return c, i.Unmarshal(c)
}

// Config for the service.
type Config struct {
	Auth           *auth.Config    `yaml:"auth,omitempty" json:"auth,omitempty" toml:"auth,omitempty"`
	Client         *client.Config  `yaml:"client,omitempty" json:"client,omitempty" toml:"client,omitempty"`
	Health         *health.Config  `yaml:"health,omitempty" json:"health,omitempty" toml:"health,omitempty"`
	Migrate        *migrate.Config `yaml:"migrate,omitempty" json:"migrate,omitempty" toml:"migrate,omitempty"`
	*config.Config `yaml:",inline" json:",inline" toml:",inline"`
}

func decorateConfig(cfg *Config) *config.Config {
	return cfg.Config
}

func v1ClientConfig(cfg *Config) *v1.Config {
	if !client.IsEnabled(cfg.Client) {
		return nil
	}

	return cfg.Client.V1
}

func v1AuthClientConfig(cfg *Config) *av1.Config {
	if !auth.IsEnabled(cfg.Auth) {
		return nil
	}

	return cfg.Auth.Client.V1
}

func healthConfig(cfg *Config) *health.Config {
	return cfg.Health
}

func migrateConfig(cfg *Config) *migrate.Config {
	return cfg.Migrate
}
