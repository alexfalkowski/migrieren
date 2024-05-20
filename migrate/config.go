package migrate

import (
	"os"
	"path/filepath"
	"strings"
)

type (
	// URL for migrate.
	URL string

	// Database for migrate.
	Database struct {
		Name   string `yaml:"name,omitempty" json:"name,omitempty" toml:"name,omitempty"`
		Source string `yaml:"source,omitempty" json:"source,omitempty" toml:"source,omitempty"`
		URL    URL    `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`
	}

	// Config for migrate.
	Config struct {
		Databases []*Database `yaml:"databases,omitempty" json:"databases,omitempty" toml:"databases,omitempty"`
	}
)

// Database by name.
func (c *Config) Database(name string) *Database {
	for _, d := range c.Databases {
		if d.Name == name {
			return d
		}
	}

	return nil
}

// GetURL for database.
func (d *Database) GetURL() (string, error) {
	k, err := os.ReadFile(filepath.Clean(string(d.URL)))

	return strings.TrimSpace(string(k)), err
}
