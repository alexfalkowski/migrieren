package migrate

import (
	"math"

	"github.com/alexfalkowski/go-service/v3/context"
	"github.com/alexfalkowski/go-service/v3/meta"
	v1 "github.com/alexfalkowski/migrieren/v2/api/migrieren/v1"
)

// Migrate migrates a configured database to the requested version.
func (s *Migrator) Migrate(ctx context.Context, req *v1.MigrateRequest) (*v1.MigrateResponse, error) {
	db := req.GetDatabase()
	ver := req.GetVersion()

	if ver == 0 || ver > math.MaxInt {
		return nil, ErrInvalidVersion
	}

	ctx, logs, err := s.migrator.Migrate(ctx, db, ver)
	if err != nil {
		return nil, err
	}

	resp := &v1.MigrateResponse{
		Migration: &v1.Migration{
			Database: db,
			Version:  ver,
			Logs:     logs,
		},
		Meta: meta.Strings(ctx),
	}

	return resp, nil
}
