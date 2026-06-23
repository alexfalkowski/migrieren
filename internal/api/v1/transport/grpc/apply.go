package grpc

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/strings"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

// ApplyMigrations applies all pending up migrations for a configured database.
func (s *Server) ApplyMigrations(ctx context.Context, req *v1.ApplyMigrationsRequest) (*v1.ApplyMigrationsResponse, error) {
	db := req.GetDatabase()
	resp := &v1.ApplyMigrationsResponse{
		Migration: &v1.Migration{Database: db},
	}

	ctx, version, logs, err := s.migrator.ApplyMigrations(ctx, db)
	setFailureTrailer(ctx, diagnostics.FromError(err))

	resp.Meta = meta.CamelStrings(ctx, strings.Empty)
	resp.Migration.Version = version
	resp.Migration.Logs = logs

	return resp, s.error(err)
}
