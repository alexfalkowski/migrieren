package config

import (
	av1 "github.com/alexfalkowski/auth/client/v1/config"
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/config"
	v1 "github.com/alexfalkowski/migrieren/client/v1/config"
	"github.com/alexfalkowski/migrieren/health"
	"github.com/alexfalkowski/migrieren/migrate"
)

// NewConfigurator for config.
func NewConfigurator(i *cmd.InputConfig) (config.Configurator, error) {
	c := &Config{}

	return c, i.Unmarshal(c)
}

func v1ClientConfig(cfg config.Configurator) *v1.Config {
	c := cfg.(*Config)
	if c.Client == nil {
		return nil
	}

	return c.Client.V1
}

func v1AuthClientConfig(cfg config.Configurator) *av1.Config {
	c := cfg.(*Config)
	if c.Auth == nil || c.Auth.Client == nil {
		return nil
	}

	return c.Auth.Client.V1
}

func healthConfig(cfg config.Configurator) *health.Config {
	return cfg.(*Config).Health
}

func migrateConfig(cfg config.Configurator) *migrate.Config {
	return cfg.(*Config).Migrate
}
