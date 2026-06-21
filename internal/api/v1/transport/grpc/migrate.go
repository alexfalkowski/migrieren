package grpc

import (
	"math"
	"strconv"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	grpcmeta "github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/api/migrate"
)

// Migrate executes the requested migration and returns response metadata and
// collected migration logs.
func (s *Server) Migrate(ctx context.Context, req *v1.MigrateRequest) (*v1.MigrateResponse, error) {
	db := req.GetDatabase()
	ver := req.GetVersion()
	resp := &v1.MigrateResponse{
		Migration: &v1.Migration{
			Database: db,
			Version:  ver,
		},
	}

	if ver == 0 || ver > math.MaxInt {
		return resp, status.Error(codes.InvalidArgument, "version must be between 1 and math.MaxInt")
	}

	ctx, logs, err := s.migrator.Migrate(ctx, db, ver)

	resp.Meta = meta.CamelStrings(ctx, strings.Empty)
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
