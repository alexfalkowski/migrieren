package migrate

import (
	"github.com/alexfalkowski/go-signal/v2"
)

// Module wires the core migrator into the DI container.
var Module = signal.Module(
	signal.Constructor(NewMigrator),
)
