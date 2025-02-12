package migrator

import (
	"context"
)

// Migrator will migrate databases.
type Migrator interface {
	// Migrate a database to a version and returning the database logs.
	Migrate(ctx context.Context, source, db string, version uint64) ([]string, error)

	// Ping the migrator.
	Ping(ctx context.Context, source, db string) error
}
