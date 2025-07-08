package migrate

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/types/slices"
)

// ErrNotFound for migrate.
var ErrNotFound = errors.New("not found")

// Config for migrate.
type Config struct {
	Databases []*Database `yaml:"databases,omitempty" json:"databases,omitempty" toml:"databases,omitempty"`
}

// Database by name.
func (c *Config) Database(name string) (*Database, error) {
	db, ok := slices.ElemFunc(c.Databases, func(d *Database) bool { return d.Name == name })
	if !ok {
		return nil, ErrNotFound
	}

	return db, nil
}

// Database for migrate.
type Database struct {
	Name   string `yaml:"name,omitempty" json:"name,omitempty" toml:"name,omitempty"`
	Source string `yaml:"source,omitempty" json:"source,omitempty" toml:"source,omitempty"`
	URL    string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`
}

// GetSource for database.
func (d *Database) GetSource(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(d.Source)
}

// GetURL for database.
func (d *Database) GetURL(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(d.URL)
}
