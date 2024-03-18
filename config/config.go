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
	Auth           *auth.Config    `yaml:"auth,omitempty" json:"auth,omitempty" toml:"auth,omitempty"`
	Client         *client.Config  `yaml:"client,omitempty" json:"client,omitempty" toml:"client,omitempty"`
	Health         *health.Config  `yaml:"health,omitempty" json:"health,omitempty" toml:"health,omitempty"`
	Migrate        *migrate.Config `yaml:"migrate,omitempty" json:"migrate,omitempty" toml:"migrate,omitempty"`
	*config.Config `yaml:",inline" json:",inline" toml:",inline"`
}
