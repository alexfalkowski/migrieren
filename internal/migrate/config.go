package migrate

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/types/slices"
)

// ErrNotFound indicates that a named database entry is not present in [Config].
var ErrNotFound = errors.New("not found")

// Config defines migration targets available to the service.
//
// Each entry in [Databases] identifies one logical database by name and carries
// the migration source and database URL references used by API handlers.
type Config struct {
	// Databases is the list of named migration targets available to callers.
	Databases []*Database `yaml:"databases,omitempty" json:"databases,omitempty" toml:"databases,omitempty"`
}

// Database returns the configured database entry for name.
//
// It returns [ErrNotFound] when no entry in [Config.Databases] matches name.
func (c *Config) Database(name string) (*Database, error) {
	db, ok := slices.ElemFunc(c.Databases, func(d *Database) bool { return d.Name == name })
	if !ok {
		return nil, ErrNotFound
	}

	return db, nil
}

// Database describes a named migration target.
//
// Source and URL are read through go-service os.FS source resolution, allowing
// direct values or indirections such as file-backed secrets.
type Database struct {
	// Name is the logical identifier used by transport requests.
	Name string `yaml:"name,omitempty" json:"name,omitempty" toml:"name,omitempty"`
	// Source is a source reference containing the migration source URL.
	Source string `yaml:"source,omitempty" json:"source,omitempty" toml:"source,omitempty"`
	// URL is a source reference containing the database connection URL.
	URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`
}

// GetSource resolves and returns the migration source URL bytes for this
// database entry.
//
// It returns an underlying source read error when Source cannot be resolved.
func (d *Database) GetSource(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(d.Source)
}

// GetURL resolves and returns the database connection URL bytes for this
// database entry.
//
// It returns an underlying source read error when URL cannot be resolved.
func (d *Database) GetURL(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(d.URL)
}
