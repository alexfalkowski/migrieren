package transport

import (
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/migrieren/transport/grpc"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	transport.Module,
	fx.Provide(grpc.UnaryServerInterceptor),
	fx.Provide(grpc.StreamServerInterceptor),
)
