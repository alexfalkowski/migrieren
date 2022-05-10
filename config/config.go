package config

import (
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/migrieren/client"
	"github.com/alexfalkowski/migrieren/health"
	"github.com/alexfalkowski/migrieren/migrate"
)

// Config for the service.
type Config struct {
	Client        client.Config  `yaml:"client"`
	Health        health.Config  `yaml:"health"`
	Migrate       migrate.Config `yaml:"migrate"`
	config.Config `yaml:",inline"`
}
