package http

import (
	"github.com/alexfalkowski/go-service/v2/context"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/api/v1/transport/grpc"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

func migrate(server *grpc.Server) func(context.Context, *v1.MigrateRequest) (*v1.MigrateResponse, error) {
	return func(ctx context.Context, req *v1.MigrateRequest) (*v1.MigrateResponse, error) {
		resp, err := server.Migrate(ctx, req)
		setFailureHeaders(ctx, diagnostics.FromError(err))

		return resp, err
	}
}
