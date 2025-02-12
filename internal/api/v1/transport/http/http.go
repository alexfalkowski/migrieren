package http

import (
	"github.com/alexfalkowski/go-service/net/http/rpc"
	"github.com/alexfalkowski/migrieren/internal/api/v1/transport/grpc"
)

// Register for HTTP.
func Register(server *grpc.Server) {
	rpc.Route("/v1/migrate", server.Migrate)
}
