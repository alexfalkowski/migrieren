package http

import (
	"github.com/alexfalkowski/go-service/v2/context"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

// PlanMigrations reports current status and pending up migration versions for a
// configured database without applying migration files.
func (s *Server) PlanMigrations(ctx context.Context, req *v1.PlanMigrationsRequest) (*v1.PlanMigrationsResponse, error) {
	resp, err := s.migrator.PlanMigrations(ctx, req)
	if err != nil {
		setFailureHeaders(ctx, diagnostics.FromError(err))
		return nil, s.error(err)
	}

	return resp, nil
}
