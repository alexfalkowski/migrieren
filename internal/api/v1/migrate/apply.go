package migrate

import (
	"github.com/alexfalkowski/go-service/v3/context"
	"github.com/alexfalkowski/go-service/v3/meta"
	v1 "github.com/alexfalkowski/migrieren/v2/api/migrieren/v1"
)

// ApplyMigrations applies all pending up migrations for a configured database.
func (s *Migrator) ApplyMigrations(ctx context.Context, req *v1.ApplyMigrationsRequest) (*v1.ApplyMigrationsResponse, error) {
	db := req.GetDatabase()

	ctx, version, logs, err := s.migrator.ApplyMigrations(ctx, db)
	if err != nil {
		return nil, err
	}

	resp := &v1.ApplyMigrationsResponse{
		Meta: meta.Strings(ctx),
		Migration: &v1.Migration{
			Database: db,
			Logs:     logs,
			Version:  version,
		},
	}

	return resp, nil
}
