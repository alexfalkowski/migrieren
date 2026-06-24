package http

import (
	"github.com/alexfalkowski/go-service/v2/context"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/api/v1/migrate"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

func migrateDatabase(migrator *migrate.Migrator) func(context.Context, *v1.MigrateRequest) (*v1.MigrateResponse, error) {
	return func(ctx context.Context, req *v1.MigrateRequest) (*v1.MigrateResponse, error) {
		resp, err := migrator.Migrate(ctx, req)
		if err != nil {
			setFailureHeaders(ctx, diagnostics.FromError(err))
			return nil, responseError(err)
		}

		return resp, nil
	}
}
