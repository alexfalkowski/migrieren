package grpc

import (
	"github.com/alexfalkowski/go-service/v2/context"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

// Migrate executes the requested migration and returns response metadata and
// collected migration logs.
func (s *Server) Migrate(ctx context.Context, req *v1.MigrateRequest) (*v1.MigrateResponse, error) {
	resp, err := s.migrator.Migrate(ctx, req)
	if err != nil {
		setFailureTrailer(ctx, diagnostics.FromError(err))
		return nil, s.error(err)
	}

	return resp, nil
}
