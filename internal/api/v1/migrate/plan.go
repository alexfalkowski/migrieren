package migrate

import (
	"math"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/strings"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// PlanMigrations reports current status and planned migration versions for a
// configured database without applying migration files.
func (s *Migrator) PlanMigrations(ctx context.Context, req *v1.PlanMigrationsRequest) (*v1.PlanMigrationsResponse, error) {
	// Keep optional-field presence: GetTargetVersion returns zero for both an
	// omitted target and an explicitly supplied, invalid zero target.
	target := req.TargetVersion
	if target != nil && (*target == 0 || *target > math.MaxInt) {
		return nil, ErrInvalidVersion
	}

	db := req.GetDatabase()

	ctx, plan, err := s.migrator.Plan(ctx, db, target)
	if err != nil {
		return nil, err
	}

	resp := &v1.PlanMigrationsResponse{
		Meta: meta.CamelStrings(ctx, strings.Empty),
		Plan: &v1.MigrationPlan{
			Status: &v1.MigrationStatus{
				Database: db,
				Version:  plan.Status.Version,
				State:    migrationState(plan.Status.State),
			},
			LatestVersion:   plan.LatestVersion,
			TargetVersion:   plan.TargetVersion,
			Direction:       migrationDirection(plan.Direction),
			PendingVersions: plan.PendingVersions,
		},
	}

	return resp, nil
}

func migrationDirection(direction migrate.PlanDirection) v1.MigrationDirection {
	switch direction {
	case migrate.PlanDirectionNone:
		return v1.MigrationDirection_MIGRATION_DIRECTION_NONE
	case migrate.PlanDirectionUp:
		return v1.MigrationDirection_MIGRATION_DIRECTION_UP
	case migrate.PlanDirectionDown:
		return v1.MigrationDirection_MIGRATION_DIRECTION_DOWN
	default:
		return v1.MigrationDirection_MIGRATION_DIRECTION_UNSPECIFIED
	}
}
