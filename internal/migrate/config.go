package migrate

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/slices"
)

// ErrNotFound is returned when a named database is not present in [Config].
var ErrNotFound = errors.New("not found")

// Config defines the configured migration targets for the service.
//
// Databases must be non-empty and each entry must have a unique Name. The
// service validates that every configured entry is present before startup
// completes.
type Config struct {
	Databases []*Database `yaml:"databases" json:"databases" toml:"databases" validate:"gt=0,unique=Name,dive,required"`
}

// Database returns the configured database entry with name.
//
// It returns [ErrNotFound] when no matching database exists.
func (c *Config) Database(name string) (*Database, error) {
	db, ok := slices.ElemFunc(c.Databases, func(d *Database) bool { return d.Name == name })
	if !ok {
		return nil, ErrNotFound
	}

	return db, nil
}

// Database describes one named migration target.
//
// Source and URL are resolved through the service filesystem abstraction so they
// can point at literal values or external sources such as `file:...`.
//
// Name, Source, and URL are all required configuration fields.
type Database struct {
	Name   string `yaml:"name,omitempty" json:"name,omitempty" toml:"name,omitempty" validate:"required"`
	Source string `yaml:"source,omitempty" json:"source,omitempty" toml:"source,omitempty" validate:"required"`
	URL    string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty" validate:"required"`
}

// GetSource resolves the configured migration source for d via fs.
func (d *Database) GetSource(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(d.Source)
}

// GetURL resolves the configured database URL for d via fs.
func (d *Database) GetURL(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(d.URL)
}
