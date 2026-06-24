package migrate

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/migrieren/internal/api/migrate"
)

// Module wires the v1 migration service into the DI container.
var Module = di.Module(
	migrate.Module,
	di.Constructor(NewMigrator),
)
