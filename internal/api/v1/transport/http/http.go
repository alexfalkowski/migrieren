package http

import (
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/api/v1/transport/grpc"
)

// Register for HTTP.
func Register(server *grpc.Server) {
	rpc.Route(v1.Service_Migrate_FullMethodName, server.Migrate)
}
