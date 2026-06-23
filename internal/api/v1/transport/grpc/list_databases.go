package grpc

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/strings"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
)

// ListDatabases reports configured logical database names.
func (s *Server) ListDatabases(ctx context.Context, _ *v1.ListDatabasesRequest) (*v1.ListDatabasesResponse, error) {
	databases := s.migrator.Databases()
	resp := &v1.ListDatabasesResponse{
		Meta:      meta.CamelStrings(ctx, strings.Empty),
		Databases: make([]*v1.Database, 0, len(databases)),
	}

	for _, name := range databases {
		resp.Databases = append(resp.Databases, &v1.Database{Name: name})
	}

	return resp, nil
}
