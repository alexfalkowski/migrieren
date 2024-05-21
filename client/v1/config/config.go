package config

import (
	"github.com/alexfalkowski/go-service/client"
)

type (
	// Config for client.
	Config struct {
		*client.Config `yaml:",inline" json:",inline" toml:",inline"`
		Migrate        *Migrate `yaml:"migrate,omitempty" json:"migrate,omitempty" toml:"migrate,omitempty"`
	}

	// Migrate the client.
	Migrate struct {
		Database string `yaml:"database,omitempty" json:"database,omitempty" toml:"database,omitempty"`
		Version  uint64 `yaml:"version,omitempty" json:"version,omitempty" toml:"version,omitempty"`
	}
)
