package grpc

import (
	"github.com/alexfalkowski/go-service/v2/context"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

// ApplyMigrations applies all pending up migrations for a configured database.
func (s *Server) ApplyMigrations(ctx context.Context, req *v1.ApplyMigrationsRequest) (*v1.ApplyMigrationsResponse, error) {
	resp, err := s.migrator.ApplyMigrations(ctx, req)
	if err != nil {
		setFailureTrailer(ctx, diagnostics.FromError(err))
		return nil, s.error(err)
	}

	return resp, nil
}
