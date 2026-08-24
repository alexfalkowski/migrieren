package http

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/api/v1/migrate"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

// Register exposes the v1 service methods through the HTTP RPC facade.
func Register(server *Server) {
	server.Route(v1.Service_Migrate_FullMethodName, server.Migrate)
	server.Route(v1.Service_ApplyMigrations_FullMethodName, server.ApplyMigrations)
	server.Route(v1.Service_PlanMigrations_FullMethodName, server.PlanMigrations)
	server.Route(v1.Service_Status_FullMethodName, server.Status)
	server.Route(v1.Service_ListDatabases_FullMethodName, server.ListDatabases)
}

// NewServer constructs an HTTP RPC transport adapter around service.
func NewServer(server *rpc.Server, service *migrate.Migrator) *Server {
	return &Server{migrator: service, Server: server}
}

// Server implements the migrieren.v1 HTTP RPC facade.
//
// It delegates migration work to the versioned API migrator and maps domain
// errors to HTTP response errors.
type Server struct {
	migrator *migrate.Migrator
	*rpc.Server
}

func setFailureHeaders(ctx context.Context, values diagnostics.Values) {
	header := meta.Response(ctx).Header()
	for key, value := range values {
		header.Set(key, value)
	}
}
