package migrate

import (
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module wires the core migrator into the DI container.
var Module = di.Module(
	di.Constructor(NewMigrator),
)
