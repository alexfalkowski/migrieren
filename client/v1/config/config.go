package config

import (
	"github.com/alexfalkowski/go-service/client"
)

// Config for client.
type Config struct {
	*client.Config `yaml:",inline" json:",inline" toml:",inline"`
	Database       string `yaml:"database,omitempty" json:"database,omitempty" toml:"database,omitempty"`
	Version        uint64 `yaml:"version,omitempty" json:"version,omitempty" toml:"version,omitempty"`
}
