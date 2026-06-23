package http

import (
	"github.com/alexfalkowski/go-service/v2/context"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/api/v1/transport/grpc"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

func apply(server *grpc.Server) func(context.Context, *v1.ApplyMigrationsRequest) (*v1.ApplyMigrationsResponse, error) {
	return func(ctx context.Context, req *v1.ApplyMigrationsRequest) (*v1.ApplyMigrationsResponse, error) {
		resp, err := server.ApplyMigrations(ctx, req)
		setFailureHeaders(ctx, diagnostics.FromError(err))

		return resp, err
	}
}
