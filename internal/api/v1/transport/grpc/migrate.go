package grpc

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/strings"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
)

// Migrate handles the v1 migration RPC.
//
// It delegates the migration request by database name and target version to the
// service migrator, then returns:
//   - the requested database/version,
//   - collected migration logs, and
//   - context metadata projected into the response meta map.
//
// Errors are translated by Server.error into gRPC status errors.
func (s *Server) Migrate(ctx context.Context, req *v1.MigrateRequest) (*v1.MigrateResponse, error) {
	db := req.GetDatabase()
	ver := req.GetVersion()
	resp := &v1.MigrateResponse{
		Migration: &v1.Migration{
			Database: db,
			Version:  ver,
		},
	}

	logs, err := s.migrator.Migrate(ctx, db, ver)

	resp.Meta = meta.CamelStrings(ctx, strings.Empty)
	resp.Migration.Logs = logs

	return resp, s.error(err)
}
