package migrate

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/migrieren/internal/migrate/database"
)

// ErrInvalidStatus is returned when migration status cannot be inspected.
var ErrInvalidStatus = errors.New("invalid status")

// StatusState identifies whether a database has no recorded migration, a clean
// recorded version, or a dirty recorded version.
type StatusState string

const (
	// StatusStateUnapplied means no migration version has been recorded.
	StatusStateUnapplied StatusState = "unapplied"

	// StatusStateClean means a non-dirty migration version is recorded.
	StatusStateClean StatusState = "clean"

	// StatusStateDirty means the recorded migration version is dirty.
	StatusStateDirty StatusState = "dirty"
)

// Status reports a database's current migration version state.
type Status struct {
	// State reports whether the migration version is unapplied, clean, or dirty.
	State StatusState

	// Version is the current clean or dirty migration version. It is zero when
	// State is StatusStateUnapplied, and also zero when State is
	// StatusStateDirty but no migration version was ever recorded (golang-migrate's
	// NilVersion-dirty recovery state).
	Version uint64
}

// Status opens databaseURL and reports the current migration version state
// without applying migration files.
//
// If no migration version has been recorded yet, the returned status has
// State=StatusStateUnapplied and Version=0.
//
// Status does not apply migration files, but the underlying migrate v4 database
// driver version path does not accept a request context. Migrieren checks ctx
// before and after the driver call, but cancellation cannot interrupt every
// upstream inspection path until migrate v5 adds context-aware driver APIs.
func (m *Migrator) Status(ctx context.Context, db string) (context.Context, *Status, error) {
	if err := ctx.Err(); err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("statusError", meta.Error(err)))
		return ctx, nil, migrationError(err)
	}

	driver, err := database.Open(ctx, db)
	if err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("statusError", meta.Error(err)))
		return ctx, nil, &invalidConfigError{err: err}
	}
	defer func() {
		_ = driver.Close()
	}()

	version, dirty, err := driver.Version()
	if err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("statusError", meta.Error(err)))
		return ctx, nil, ErrInvalidStatus
	}

	if err := ctx.Err(); err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair("statusError", meta.Error(err)))
		return ctx, nil, migrationError(err)
	}

	status := &Status{State: StatusStateUnapplied}
	switch {
	case dirty:
		status.State = StatusStateDirty
		if version >= 0 {
			status.Version = uint64(version)
		}
	case version >= 0:
		status.Version = uint64(version)
		status.State = StatusStateClean
	}

	return ctx, status, nil
}
