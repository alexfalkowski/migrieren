package migrate

import (
	"os"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/migrieren/internal/migrate/source"
	migratesource "github.com/golang-migrate/migrate/v4/source"
)

// PlanDirection identifies whether applying pending up migrations would advance
// the database from its current migration status.
type PlanDirection uint8

const (
	// PlanDirectionNone means no up migrations are pending.
	PlanDirectionNone PlanDirection = iota + 1

	// PlanDirectionUp means one or more up migrations are pending.
	PlanDirectionUp
)

// Plan reports non-mutating migration planning information.
type Plan struct {
	// Status reports the current database migration version state.
	Status *Status

	// PendingVersions contains source versions greater than Status.Version in
	// apply order.
	PendingVersions []uint64

	// Direction reports whether up migrations are pending.
	Direction PlanDirection

	// LatestVersion is the highest migration version available from the source.
	LatestVersion uint64

	// TargetVersion is the version a clean or unapplied database can converge
	// toward by applying pending up migrations. It equals Status.Version when no
	// up migrations can apply.
	TargetVersion uint64
}

// Plan reports current database status and available up migration versions
// without applying migration files.
//
// Inputs:
//   - ctx: a service context used for metadata/telemetry.
//   - src: migration source URL (for example "file://...").
//   - db: database URL (for example a Postgres URL).
//
// Output:
//   - ctx: the input context, or a derived context containing "planError" or
//     "statusError" when source/database setup or inspection fails.
//   - plan: current status, latest source version, target version, direction,
//     and pending source versions.
//   - error: nil on success; otherwise one of:
//     [ErrInvalidConfig] (cannot open src/db),
//     [ErrInvalidStatus] (database status cannot be inspected),
//     [ErrInvalidMigration] (source versions cannot be inspected),
//     [ErrMigrationCanceled] (request context canceled),
//     [ErrMigrationDeadlineExceeded] (request context deadline expired).
//
// Plan opens the migration source to inspect available versions, but does not
// read migration bodies or execute migration operations.
func (m *Migrator) Plan(ctx context.Context, src, db string) (context.Context, *Plan, error) {
	if err := ctx.Err(); err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("planError", meta.Error(err)))
		return ctx, nil, migrationError(err)
	}

	driver, err := source.Open(src)
	if err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("planError", meta.Error(err)))
		return ctx, nil, ErrInvalidConfig
	}
	defer func() {
		_ = driver.Close()
	}()

	ctx, status, err := m.Status(ctx, db)
	if err != nil {
		return ctx, nil, err
	}

	versions, err := sourceVersions(driver)
	if err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("planError", meta.Error(err)))
		return ctx, nil, ErrInvalidMigration
	}

	if err := ctx.Err(); err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("planError", meta.Error(err)))
		return ctx, nil, migrationError(err)
	}

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

	return ctx, plan, nil
}

func sourceVersions(driver migratesource.Driver) ([]uint64, error) {
	version, err := driver.First()
	versions := make([]uint64, 0)
	for {
		if errors.Is(err, os.ErrNotExist) {
			return versions, nil
		}
		if err != nil {
			return nil, err
		}

		versions = append(versions, uint64(version))

		version, err = driver.Next(version)
	}
}
