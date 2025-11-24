package migrate

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/migrieren/internal/migrate/database"
)

// Module for fx.
var Module = di.Module(
	di.Register(database.Register),
	di.Constructor(NewMigrator),
)
