package config

import (
	"github.com/alexfalkowski/go-service/client"
)

// Config for client.
type Config struct {
	Database      string `yaml:"database" json:"database" toml:"database"`
	Version       uint64 `yaml:"version" json:"version" toml:"version"`
	client.Config `yaml:",inline" json:",inline" toml:",inline"`
}
