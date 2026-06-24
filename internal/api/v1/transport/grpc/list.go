package grpc

import (
	"github.com/alexfalkowski/go-service/v2/context"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
)

// ListDatabases reports configured logical database names.
func (s *Server) ListDatabases(ctx context.Context, req *v1.ListDatabasesRequest) (*v1.ListDatabasesResponse, error) {
	return s.migrator.ListDatabases(ctx, req)
}
