package v1

import (
	"github.com/alexfalkowski/migrieren/migrate"
	"github.com/alexfalkowski/migrieren/server/v1/transport/grpc"
	"github.com/alexfalkowski/migrieren/server/v1/transport/grpc/security/token"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(migrate.NewMigrator),
	fx.Provide(token.NewVerifier),
	fx.Provide(grpc.NewServer),
	fx.Invoke(grpc.Register),
)
