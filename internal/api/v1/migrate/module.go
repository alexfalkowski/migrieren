package migrate

import (
	"github.com/alexfalkowski/go-signal/v2"
	"github.com/alexfalkowski/migrieren/internal/api/migrate"
)

// Module wires the v1 migration service into the DI container.
var Module = signal.Module(
	migrate.Module,
	signal.Constructor(NewMigrator),
)
