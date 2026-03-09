package v1

import (
	"github.com/alexfalkowski/go-service/v2/di"
	api "github.com/alexfalkowski/migrieren/internal/api/migrate"
	"github.com/alexfalkowski/migrieren/internal/api/v1/transport/grpc"
	"github.com/alexfalkowski/migrieren/internal/api/v1/transport/http"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// Module registers API v1 dependencies for DI.
//
// It composes migration modules, constructs the gRPC server, and registers both
// gRPC and HTTP transport handlers.
var Module = di.Module(
	api.Module,
	migrate.Module,
	di.Constructor(grpc.NewServer),
	di.Register(grpc.Register),
	di.Register(http.Register),
)
