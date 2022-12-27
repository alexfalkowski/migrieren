package v1

import (
	"github.com/alexfalkowski/migrieren/migrate"
	"github.com/alexfalkowski/migrieren/migrate/trace/opentracing"
	"github.com/alexfalkowski/migrieren/server/v1/transport/grpc"
	"go.uber.org/fx"
)

var (
	// Module for fx.
	Module = fx.Options(
		fx.Provide(migrate.NewMigrator),
		fx.Provide(opentracing.NewTracer),
		fx.Provide(grpc.NewServer),
		fx.Invoke(grpc.Register),
	)
)
