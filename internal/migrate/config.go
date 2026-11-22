package migrate

import (
	"github.com/alexfalkowski/go-service/v3/database/sql/config"
	"github.com/alexfalkowski/go-service/v3/errors"
	"github.com/alexfalkowski/go-service/v3/os"
	"github.com/alexfalkowski/go-service/v3/slices"
)

// ErrNotFound is returned when a named database is not present in [Config].
var ErrNotFound = errors.New("not found")

// Config defines the configured migration targets for the service.
//
// Databases must be non-empty and each entry must have a unique Name. The
// service validates that every configured entry is present before startup
// completes.
//
// Logs is optional and tunes how many migration log lines each operation
// returns. When Logs is nil or its Max is not positive, the service uses its
// default maximum.
type Config struct {
	Logs      *Logs       `yaml:"logs,omitempty" json:"logs,omitempty" toml:"logs,omitempty" validate:"omitempty"`
	Databases []*Database `yaml:"databases" json:"databases" toml:"databases" validate:"gt=0,unique=Name,dive,required"`
}

// defaultMaxLogs is the number of migration log lines returned when a database
// does not configure a positive migrate.logs.max.
const defaultMaxLogs = 100

// Logs configures how migration log output is captured and returned.
type Logs struct {
	// Max is the maximum number of migration log lines returned per operation.
	// When not positive, the default maximum is used. When exceeded, the oldest
	// lines are dropped and the first returned entry is a truncation marker.
	Max int `yaml:"max,omitempty" json:"max,omitempty" toml:"max,omitempty" validate:"omitempty,gte=0"`
}

// GetMax returns the configured maximum number of migration log lines, or the
// default maximum when Logs is nil or Max is not positive.
//
// It is safe to call on a nil receiver, so callers can use cfg.Logs.GetMax()
// without first checking whether the optional logs section was configured.
func (l *Logs) GetMax() int {
	if l == nil || l.Max <= 0 {
		return defaultMaxLogs
	}

	return l.Max
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
// A target requires a name, a migration source, and a database connection source.
type Database struct {
	// Name is the logical database name used by callers.
	Name string `yaml:"name,omitempty" json:"name,omitempty" toml:"name,omitempty" validate:"required"`
	// Source resolves to the migration files' URL.
	Source string `yaml:"source,omitempty" json:"source,omitempty" toml:"source,omitempty" validate:"required"`
	// DSN resolves to the database connection URL.
	DSN config.DSN `yaml:"dsn,omitempty" json:"dsn,omitzero" toml:"dsn,omitempty" validate:"required"`
}

// GetSource resolves the configured migration source for d via fs.
func (d *Database) GetSource(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(d.Source)
}

// GetURL resolves the database connection URL from the DSN source and returns
// any error from fs.
func (d *Database) GetURL(fs *os.FS) ([]byte, error) {
	return d.DSN.Resolve(fs)
}
