package migrate

import (
	"slices"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/migrieren/internal/migrate/source"
)

// PlanDirection identifies whether a migration plan advances, rolls back, or
// preserves the database's current migration status.
type PlanDirection uint8

const (
	// PlanDirectionNone means no migrations are pending.
	PlanDirectionNone PlanDirection = iota + 1

	// PlanDirectionUp means one or more up migrations are pending.
	PlanDirectionUp

	// PlanDirectionDown means one or more down migrations are pending.
	PlanDirectionDown
)

// Plan reports non-mutating migration planning information.
type Plan struct {
	// Status reports the current database migration version state.
	Status *Status

	// PendingVersions contains source versions that the plan would apply or
	// remove, in execution order.
	PendingVersions []uint64

	// Direction reports whether up or down migrations are pending.
	Direction PlanDirection

	// LatestVersion is the highest migration version available from the source.
	LatestVersion uint64

	// TargetVersion is the version the plan can converge toward.
	TargetVersion uint64
}

// Plan reports current database status and available migration versions without
// applying migration files.
//
// Inputs:
//   - ctx: a service context used for metadata/telemetry.
//   - src: migration source URL (for example "file://...").
//   - db: database URL (for example a Postgres URL).
//   - target: optional migration version to preview. When nil, Plan preserves
//     the latest-up planning behavior.
//
// Output:
//   - ctx: the input context, or a derived context containing "planError" or
//     "statusError" when source/database setup or inspection fails.
//   - plan: current status, latest source version, target version, direction,
//     and pending source versions.
//   - error: nil on success; otherwise one of:
//     [ErrInvalidConfig] (cannot open src/db),
//     [ErrInvalidStatus] (database status cannot be inspected),
//     [ErrInvalidMigration] (source versions cannot be inspected or the
//     requested target cannot be traversed),
//     [ErrMigrationCanceled] (request context canceled),
//     [ErrMigrationDeadlineExceeded] (request context deadline expired).
//
// Plan opens the migration source to inspect available versions, but does not
// read migration bodies or execute migration operations. A supplied target, and
// a non-unapplied current version, must be present in the source; otherwise Plan
// returns [ErrInvalidMigration].
func (m *Migrator) Plan(ctx context.Context, src, db string, target *uint64) (context.Context, *Plan, error) {
	if err := ctx.Err(); err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("planError", meta.Error(err)))
		return ctx, nil, migrationError(err)
	}

	driver, err := source.Open(src)
	if err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("planError", meta.Error(err)))
		return ctx, nil, &invalidConfigError{err: err}
	}
	defer func() {
		_ = driver.Close()
	}()

	ctx, status, err := m.Status(ctx, db)
	if err != nil {
		return ctx, nil, err
	}

	versions, err := source.Versions(driver)
	if err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("planError", meta.Error(err)))
		return ctx, nil, ErrInvalidMigration
	}

	if err := ctx.Err(); err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("planError", meta.Error(err)))
		return ctx, nil, migrationError(err)
	}

	if target == nil {
		return ctx, latestPlan(status, versions), nil
	}

	return planForTarget(ctx, status, versions, *target)
}

func latestPlan(status *Status, versions []uint64) *Plan {
	plan := &Plan{
		Status:        status,
		TargetVersion: status.Version,
		Direction:     PlanDirectionNone,
	}
	for _, version := range versions {
		plan.LatestVersion = version
		if version > status.Version {
			plan.PendingVersions = append(plan.PendingVersions, version)
		}
	}

	if len(plan.PendingVersions) > 0 && status.State != StatusStateDirty {
		plan.Direction = PlanDirectionUp
		plan.TargetVersion = plan.LatestVersion
	}

	return plan
}

func planForTarget(ctx context.Context, status *Status, versions []uint64, target uint64) (context.Context, *Plan, error) {
	if status.State == StatusStateDirty {
		return invalidPlan(ctx)
	}

	plan := targetPlan(status, versions, target)
	if plan == nil {
		return invalidPlan(ctx)
	}

	return ctx, plan, nil
}

func targetPlan(status *Status, versions []uint64, target uint64) *Plan {
	plan := &Plan{
		Status:        status,
		TargetVersion: target,
		Direction:     PlanDirectionNone,
	}
	currentFound := status.State == StatusStateUnapplied
	targetFound := false
	for _, version := range versions {
		plan.LatestVersion = version
		currentFound = currentFound || version == status.Version
		targetFound = targetFound || version == target
	}
	if !currentFound || !targetFound {
		return nil
	}

	plan.Direction, plan.PendingVersions = targetPendingVersions(status.Version, versions, target)

	return plan
}

func targetPendingVersions(current uint64, versions []uint64, target uint64) (PlanDirection, []uint64) {
	switch {
	case current < target:
		return PlanDirectionUp, upPendingVersions(current, versions, target)
	case current > target:
		return PlanDirectionDown, downPendingVersions(current, versions, target)
	default:
		return PlanDirectionNone, nil
	}
}

func upPendingVersions(current uint64, versions []uint64, target uint64) []uint64 {
	pending := make([]uint64, 0)
	for _, version := range versions {
		if version > current && version <= target {
			pending = append(pending, version)
		}
	}

	return pending
}

func downPendingVersions(current uint64, versions []uint64, target uint64) []uint64 {
	pending := make([]uint64, 0)
	for _, version := range slices.Backward(versions) {
		if version <= target {
			break
		}
		if version <= current {
			pending = append(pending, version)
		}
	}

	return pending
}

func invalidPlan(ctx context.Context) (context.Context, *Plan, error) {
	ctx = meta.WithAttributes(ctx, meta.NewPair("planError", meta.Error(ErrInvalidMigration)))
	return ctx, nil, ErrInvalidMigration
}
