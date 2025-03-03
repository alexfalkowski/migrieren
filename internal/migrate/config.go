package migrate

import (
	"github.com/alexfalkowski/go-service/os"
)

// Config for migrate.
type Config struct {
	Databases []*Database `yaml:"databases,omitempty" json:"databases,omitempty" toml:"databases,omitempty"`
}

// Database by name.
func (c *Config) Database(name string) *Database {
	for _, d := range c.Databases {
		if d.Name == name {
			return d
		}
	}

	return nil
}

// Database for migrate.
type Database struct {
	Name   string `yaml:"name,omitempty" json:"name,omitempty" toml:"name,omitempty"`
	Source string `yaml:"source,omitempty" json:"source,omitempty" toml:"source,omitempty"`
	URL    URL    `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`
}

// GetURL for database.
func (d *Database) GetURL() ([]byte, error) {
	return os.ReadFile(d.URL.String())
}

// URL for migrate.
type URL string

// String for URL.
func (u URL) String() string {
	return string(u)
}
