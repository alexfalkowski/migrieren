package migrate

import (
	"github.com/alexfalkowski/go-signal/v2"
	"github.com/alexfalkowski/migrieren/v2/internal/migrate"
)

// Module wires the API-facing migrator into the DI container.
var Module = signal.Module(
	migrate.Module,
	signal.Constructor(NewMigrator),
)
