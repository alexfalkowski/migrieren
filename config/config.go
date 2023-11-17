package config

import (
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/migrieren/auth"
	"github.com/alexfalkowski/migrieren/client"
	"github.com/alexfalkowski/migrieren/health"
	"github.com/alexfalkowski/migrieren/migrate"
)

// Config for the service.
type Config struct {
	Auth          auth.Config    `yaml:"auth" json:"auth" toml:"auth"`
	Client        client.Config  `yaml:"client" json:"client" toml:"client"`
	Health        health.Config  `yaml:"health" json:"health" toml:"health"`
	Migrate       migrate.Config `yaml:"migrate" json:"migrate" toml:"migrate"`
	config.Config `yaml:",inline" json:",inline" toml:",inline"`
}
