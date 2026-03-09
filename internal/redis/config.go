package redis

// Config defines Redis client settings used by this package.
type Config struct {
	// URL points to a Redis connection URL source.
	//
	// The value is read through go-service os.FS source resolution, so it can be
	// a direct URL or an indirection such as "file:...".
	URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`
}
