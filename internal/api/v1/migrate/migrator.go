package migrate

import "github.com/alexfalkowski/migrieren/internal/api/migrate"

// NewMigrator constructs a v1 migration service.
//
// Dependencies:
//   - migrator resolves configured database names and executes migration work.
func NewMigrator(migrator *migrate.Migrator) *Migrator {
	return &Migrator{migrator: migrator}
}

// Migrator implements the transport-neutral migrieren.v1 migration contract.
//
// It accepts generated v1 request messages, delegates database-name migration
// work to the API migrator, and returns generated v1 response messages for
// transport adapters to expose.
type Migrator struct {
	migrator *migrate.Migrator
}
