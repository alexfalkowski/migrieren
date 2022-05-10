package client

import (
	"time"
)

// Config for client.
type Config struct {
	Host     string        `yaml:"host"`
	Timeout  time.Duration `yaml:"timeout"`
	Database string        `yaml:"database"`
	Version  uint64        `yaml:"version"`
}
