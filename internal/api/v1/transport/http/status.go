package http

import (
	"github.com/alexfalkowski/go-service/v2/context"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/api/v1/migrate"
)

func getStatus(migrator *migrate.Migrator) func(context.Context, *v1.StatusRequest) (*v1.StatusResponse, error) {
	return func(ctx context.Context, req *v1.StatusRequest) (*v1.StatusResponse, error) {
		resp, err := migrator.Status(ctx, req)
		if err != nil {
			return nil, responseError(err)
		}

		return resp, nil
	}
}
