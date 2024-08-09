package v1

import (
	"github.com/alexfalkowski/migrieren/migrate"
	sm "github.com/alexfalkowski/migrieren/server/migrate"
	"github.com/alexfalkowski/migrieren/server/v1/transport/grpc"
	"github.com/alexfalkowski/migrieren/server/v1/transport/http"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	sm.Module,
	migrate.Module,
	fx.Provide(grpc.NewServer),
	fx.Invoke(grpc.Register),
	fx.Invoke(http.Register),
)
