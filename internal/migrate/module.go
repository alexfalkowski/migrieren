package migrate

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/migrieren/internal/redis"
)

// Module registers migration dependencies for DI.
//
// It includes [redis.Module] and provides [NewMigrator].
var Module = di.Module(
	redis.Module,
	di.Constructor(NewMigrator),
)
