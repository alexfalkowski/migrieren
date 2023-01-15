package config

import (
	"time"
)

// Config for client.
type Config struct {
	Host     string        `yaml:"host" json:"host" toml:"host"`
	Timeout  time.Duration `yaml:"timeout" json:"timeout" toml:"timeout"`
	Database string        `yaml:"database" json:"database" toml:"database"`
	Version  uint64        `yaml:"version" json:"version" toml:"version"`
}
