package health

// Config defines health scheduler and timeout settings.
type Config struct {
	// Duration is the check interval (for example "30s").
	Duration string `yaml:"duration,omitempty" json:"duration,omitempty" toml:"duration,omitempty"`
	// Timeout is the per-check timeout budget (for example "5s").
	Timeout string `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`
}
