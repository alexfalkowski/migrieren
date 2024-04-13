package v1

import (
	"github.com/alexfalkowski/migrieren/migrate"
	"github.com/alexfalkowski/migrieren/server/v1/transport/grpc"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(migrate.NewMigrator),
	fx.Provide(grpc.NewServer),
	fx.Invoke(grpc.Register),
)
