package http

import (
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/api/v1/transport/grpc"
)

// Register registers HTTP RPC routes for the v1 service.
//
// The Migrate route is forwarded to the gRPC server implementation.
func Register(server *grpc.Server) {
	rpc.Route(v1.Service_Migrate_FullMethodName, server.Migrate)
}
