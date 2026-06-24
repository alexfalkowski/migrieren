package migrate

import (
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// NewMigrator constructs an API-facing [Migrator].
//
// Dependencies:
//   - migrator executes core migration operations against resolved source and
//     database URLs.
//   - fs resolves configured source and URL values.
//   - cfg contains the configured logical database list.
func NewMigrator(migrator *migrate.Migrator, fs *os.FS, cfg *migrate.Config) *Migrator {
	return &Migrator{migrator: migrator, fs: fs, config: cfg}
}

// Migrator adapts the core migrator to a database-name API.
//
// It resolves configured source and URL values for a logical database name,
// delegates work to the core migrator, and keeps configured secrets out of the
// versioned API layer.
type Migrator struct {
	migrator *migrate.Migrator
	config   *migrate.Config
	fs       *os.FS
}
