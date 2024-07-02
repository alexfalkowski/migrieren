package grpc

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
)

// Migrate for gRPC.
func (s *Server) Migrate(ctx context.Context, req *v1.MigrateRequest) (*v1.MigrateResponse, error) {
	db := req.GetDatabase()
	ver := req.GetVersion()
	resp := &v1.MigrateResponse{
		Migration: &v1.Migration{
			Database: db,
			Version:  ver,
		},
	}

	logs, err := s.service.Migrate(ctx, db, ver)
	if err != nil {
		resp.Meta = meta.CamelStrings(ctx, "")

		return resp, s.error(err)
	}

	resp.Meta = meta.CamelStrings(ctx, "")
	resp.Migration.Logs = logs

	return resp, nil
}
