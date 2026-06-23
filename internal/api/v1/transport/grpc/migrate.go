package grpc

import (
	"math"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	grpcmeta "github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
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
	setFailureTrailer(ctx, diagnostics.FromError(err))

	resp.Meta = meta.CamelStrings(ctx, strings.Empty)
	resp.Migration.Logs = logs

	return resp, s.error(err)
}

func setFailureTrailer(ctx context.Context, values diagnostics.Values) {
	pairs := failureDiagnosticPairs(values)
	if len(pairs) == 0 {
		return
	}

	_ = grpc.SetTrailer(ctx, grpcmeta.Pairs(pairs...))
}

func failureDiagnosticPairs(values diagnostics.Values) []string {
	pairs := make([]string, 0, len(values)*2)
	for key, value := range values {
		pairs = append(pairs, key, value)
	}

	return pairs
}
