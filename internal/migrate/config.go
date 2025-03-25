package migrate

import (
	"slices"

	"github.com/alexfalkowski/go-service/os"
)

// Config for migrate.
type Config struct {
	Databases []*Database `yaml:"databases,omitempty" json:"databases,omitempty" toml:"databases,omitempty"`
}

// Database by name.
func (c *Config) Database(name string) *Database {
	index := slices.IndexFunc(c.Databases, func(d *Database) bool { return d.Name == name })
	if index == -1 {
		return nil
	}

	return c.Databases[index]
}

// Database for migrate.
type Database struct {
	Name   string `yaml:"name,omitempty" json:"name,omitempty" toml:"name,omitempty"`
	Source string `yaml:"source,omitempty" json:"source,omitempty" toml:"source,omitempty"`
	URL    string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`
}

// GetSource for database.
func (d *Database) GetSource() ([]byte, error) {
	return os.ReadFile(d.Source)
}

// GetURL for database.
func (d *Database) GetURL() ([]byte, error) {
	return os.ReadFile(d.URL)
}
