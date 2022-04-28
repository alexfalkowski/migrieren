package config

import (
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/migrieren/health"
)

// Config for the service.
type Config struct {
	Health        health.Config `yaml:"health"`
	config.Config `yaml:",inline"`
}
