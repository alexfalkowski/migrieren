package config

import (
	"github.com/alexfalkowski/go-service/v2/config"
	"github.com/alexfalkowski/migrieren/internal/health"
	"github.com/alexfalkowski/migrieren/internal/migrate"
	"github.com/alexfalkowski/migrieren/internal/redis"
)

// Config is the root service configuration model.
//
// It includes subsystem configuration blocks plus the shared base
// go-service configuration via embedded [config.Config].
type Config struct {
	// Health configures service health checks and endpoints.
	Health *health.Config `yaml:"health,omitempty" json:"health,omitempty" toml:"health,omitempty"`
	// Migrate configures named migration database targets.
	Migrate *migrate.Config `yaml:"migrate,omitempty" json:"migrate,omitempty" toml:"migrate,omitempty"`
	// Redis configures the distributed lock client used by migration operations.
	Redis *redis.Config `yaml:"redis,omitempty" json:"redis,omitempty" toml:"redis,omitempty"`
	// Config is the shared go-service base configuration.
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

func redisConfig(cfg *Config) *redis.Config {
	return cfg.Redis
}
