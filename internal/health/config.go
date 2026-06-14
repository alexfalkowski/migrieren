package health

import "github.com/alexfalkowski/go-service/v2/time"

// Config defines health check timing configuration.
//
// Duration and Timeout are duration strings (for example "5s", "1m") that are
// parsed as Go durations.
//
// Duration controls how often health registrations are evaluated/updated.
// Timeout caps online checks and each per-database source/database health
// check.
type Config struct {
	Duration time.Duration `yaml:"duration,omitempty" json:"duration,omitempty" toml:"duration,omitempty" validate:"gt=0"`
	Timeout  time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty" validate:"gt=0"`
}
