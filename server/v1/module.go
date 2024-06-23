package v1

import (
	"github.com/alexfalkowski/migrieren/migrate"
	"github.com/alexfalkowski/migrieren/server/service"
	"github.com/alexfalkowski/migrieren/server/v1/transport/grpc"
	"github.com/alexfalkowski/migrieren/server/v1/transport/grpc/security/token"
	"github.com/alexfalkowski/migrieren/server/v1/transport/http"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	service.Module,
	migrate.Module,
	fx.Provide(token.NewVerifier),
	fx.Provide(grpc.NewServer),
	fx.Invoke(grpc.Register),
	fx.Invoke(http.Register),
)
