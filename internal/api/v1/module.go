package v1

import (
	api "github.com/alexfalkowski/migrieren/internal/api/migrate"
	"github.com/alexfalkowski/migrieren/internal/api/v1/transport/grpc"
	"github.com/alexfalkowski/migrieren/internal/api/v1/transport/http"
	"github.com/alexfalkowski/migrieren/internal/migrate"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	api.Module,
	migrate.Module,
	fx.Provide(grpc.NewServer),
	fx.Invoke(grpc.Register),
	fx.Invoke(http.Register),
)
