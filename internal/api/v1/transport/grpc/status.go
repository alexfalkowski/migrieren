package grpc

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/strings"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// Status reports the current migration version state for a configured database.
func (s *Server) Status(ctx context.Context, req *v1.StatusRequest) (*v1.StatusResponse, error) {
	db := req.GetDatabase()
	resp := &v1.StatusResponse{
		Status: &v1.MigrationStatus{Database: db},
	}

	ctx, state, err := s.migrator.Status(ctx, db)
	resp.Meta = meta.CamelStrings(ctx, strings.Empty)

	if state != nil {
		resp.Status.Version = state.Version
		resp.Status.State = migrationState(state.State)
	}

	return resp, s.error(err)
}

func migrationState(state migrate.StatusState) v1.MigrationState {
	switch state {
	case migrate.StatusStateUnapplied:
		return v1.MigrationState_MIGRATION_STATE_UNAPPLIED
	case migrate.StatusStateClean:
		return v1.MigrationState_MIGRATION_STATE_CLEAN
	case migrate.StatusStateDirty:
		return v1.MigrationState_MIGRATION_STATE_DIRTY
	default:
		return v1.MigrationState_MIGRATION_STATE_UNSPECIFIED
	}
}
