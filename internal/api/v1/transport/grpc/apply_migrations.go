package grpc

import (
	"strconv"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	grpcmeta "github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/strings"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/api/migrate"
)

// ApplyMigrations applies all pending up migrations for a configured database.
func (s *Server) ApplyMigrations(
	ctx context.Context, req *v1.ApplyMigrationsRequest,
) (*v1.ApplyMigrationsResponse, error) {
	db := req.GetDatabase()
	resp := &v1.ApplyMigrationsResponse{
		Migration: &v1.Migration{Database: db},
	}

	ctx, version, logs, err := s.migrator.ApplyMigrations(ctx, db)

	resp.Meta = meta.CamelStrings(ctx, strings.Empty)
	resp.Migration.Version = version
	resp.Migration.Logs = logs

	if err != nil {
		values := []string{
			"migration-error", migrate.FailureKind(ctx, err),
			"migration-log-count", strconv.Itoa(len(logs)),
		}

		if stage := meta.Attribute(ctx, migrate.FailureStageKey); !stage.IsEmpty() {
			values = append(values, "migration-stage", stage.String())
		}

		if len(logs) > 0 {
			values = append(values, "migration-log-last", logs[len(logs)-1])
		}

		_ = grpc.SetTrailer(ctx, grpcmeta.Pairs(values...))
	}

	return resp, s.error(err)
}
