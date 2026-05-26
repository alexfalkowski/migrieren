package health

import "github.com/alexfalkowski/go-service/v2/time"

const (
	// DefaultDuration is used when health.duration is omitted or non-positive.
	DefaultDuration = time.Second

	// DefaultTimeout is used when health.timeout is omitted or non-positive.
	DefaultTimeout = time.Second
)

// Config defines health check timing configuration.
//
// Duration and Timeout are duration strings (for example "5s", "1m") that are
// parsed as Go durations.
//
// Duration controls how often health registrations are evaluated/updated.
// Timeout is reserved for probe/check timeout configuration.
type Config struct {
	Duration time.Duration `yaml:"duration,omitempty" json:"duration,omitempty" toml:"duration,omitempty"`
	Timeout  time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`
}

// GetDuration returns the configured health check duration or a default.
func (c *Config) GetDuration() time.Duration {
	if c.Duration > 0 {
		return c.Duration
	}

	return DefaultDuration
}

// GetTimeout returns the configured health check timeout or a default.
func (c *Config) GetTimeout() time.Duration {
	if c.Timeout > 0 {
		return c.Timeout
	}

	return DefaultTimeout
}
