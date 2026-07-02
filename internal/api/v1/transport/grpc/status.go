package grpc

import (
	"github.com/alexfalkowski/go-service/v2/context"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

// Status reports the current migration version state for a configured database.
func (s *Server) Status(ctx context.Context, req *v1.StatusRequest) (*v1.StatusResponse, error) {
	resp, err := s.migrator.Status(ctx, req)
	if err != nil {
		setFailureTrailer(ctx, diagnostics.FromError(err))
		return nil, s.error(err)
	}

	return resp, nil
}
