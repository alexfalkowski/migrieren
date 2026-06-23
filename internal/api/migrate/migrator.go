package migrate

import (
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// NewMigrator constructs a transport-facing [Migrator].
//
// Dependencies:
//   - migrator: the core migrator that executes migrations given a source URL and
//     database URL.
//   - fs: a filesystem abstraction used to resolve `Database.Source` and
//     `Database.URL` values (for example `file:...`).
//   - cfg: the migration configuration containing the list of named databases.
func NewMigrator(migrator *migrate.Migrator, fs *os.FS, cfg *migrate.Config) *Migrator {
	return &Migrator{migrator: migrator, fs: fs, config: cfg}
}

// Migrator adapts the core migrator to a "database name + version" API that is
// convenient for transport layers.
//
// The adapter:
//   - looks up a database by name in the provided config,
//   - reads its source and URL through the filesystem abstraction,
//   - delegates migration execution or status inspection to the core migrator,
//   - lists configured logical database names without exposing source or URL values.
type Migrator struct {
	migrator *migrate.Migrator
	config   *migrate.Config
	fs       *os.FS
}
