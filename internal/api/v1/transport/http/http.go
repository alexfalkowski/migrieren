package http

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/api/v1/transport/grpc"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

// Register exposes the gRPC service methods through the HTTP RPC facade.
func Register(server *grpc.Server) {
	rpc.Route(v1.Service_Migrate_FullMethodName, migrate(server))
	rpc.Route(v1.Service_ApplyMigrations_FullMethodName, apply(server))
	rpc.Route(v1.Service_Status_FullMethodName, server.Status)
	rpc.Route(v1.Service_ListDatabases_FullMethodName, server.ListDatabases)
}

func setFailureHeaders(ctx context.Context, values diagnostics.Values) {
	header := meta.Response(ctx).Header()
	for key, value := range values {
		header.Set(key, value)
	}
}
