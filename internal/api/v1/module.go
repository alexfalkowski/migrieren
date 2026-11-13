package v1

import (
	"github.com/alexfalkowski/go-signal/v2"
	"github.com/alexfalkowski/migrieren/internal/api/v1/migrate"
	"github.com/alexfalkowski/migrieren/internal/api/v1/transport/grpc"
	"github.com/alexfalkowski/migrieren/internal/api/v1/transport/http"
)

// Module wires the v1 API transports and handlers into the DI container.
var Module = signal.Module(
	migrate.Module,
	signal.Constructor(grpc.NewServer),
	signal.Constructor(http.NewServer),
	signal.Register(grpc.Register),
	signal.Register(http.Register),
)
