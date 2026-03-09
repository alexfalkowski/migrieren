package migrate

import "github.com/alexfalkowski/go-service/v2/di"

// Module registers the transport-facing migration adapter for DI.
//
// It provides [NewMigrator].
var Module = di.Module(
	di.Constructor(NewMigrator),
)
