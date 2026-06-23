package diagnostics

import (
	"maps"
	"strconv"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// Values contains safe diagnostic key/value pairs attached to an error.
type Values map[string]string

const (
	stageKey = "migrateErrorStage"

	// StageSource identifies failures while resolving a migration source.
	StageSource = "source"

	// StageURL identifies failures while resolving a database URL.
	StageURL = "url"
)

// WithStage stores the migration diagnostic stage on ctx.
func WithStage(ctx context.Context, stage string) context.Context {
	return meta.WithAttributes(ctx, meta.NewPair(stageKey, meta.String(stage)))
}

// Error wraps err with safe migration diagnostic values.
//
// The returned error unwraps to err, so errors.Is and errors.As continue to
// match the original cause. Passing nil returns nil.
func Error(ctx context.Context, err error, logs []string) error {
	if err == nil {
		return nil
	}

	return &diagnosticError{
		err:    err,
		values: newValues(ctx, err, logs),
	}
}

func newValues(ctx context.Context, err error, logs []string) Values {
	values := Values{
		"migration-error":     failureKind(ctx, err),
		"migration-log-count": strconv.Itoa(len(logs)),
	}

	if stage := meta.Attribute(ctx, stageKey); !stage.IsEmpty() {
		values["migration-stage"] = stage.String()
	}

	if len(logs) > 0 {
		values["migration-log-last"] = logs[len(logs)-1]
	}

	return values
}

func failureKind(ctx context.Context, err error) string {
	switch {
	case errors.Is(err, migrate.ErrNotFound):
		return "not_found"
	case errors.Is(err, migrate.ErrMigrationCanceled):
		return "canceled"
	case errors.Is(err, migrate.ErrMigrationDeadlineExceeded):
		return "deadline_exceeded"
	case !meta.Attribute(ctx, stageKey).IsEmpty(), errors.Is(err, migrate.ErrInvalidConfig):
		return "invalid_config"
	case errors.Is(err, migrate.ErrInvalidMigration):
		return "invalid_migration"
	default:
		return "unknown"
	}
}

// FromError returns the safe diagnostic values carried by err.
//
// If err does not carry diagnostics, FromError returns an empty map.
func FromError(err error) Values {
	if diagnostic, ok := errors.AsType[*diagnosticError](err); ok {
		return diagnostic.values.copy()
	}

	return Values{}
}

type diagnosticError struct {
	err    error
	values Values
}

func (d *diagnosticError) Error() string {
	return d.err.Error()
}

func (d *diagnosticError) Unwrap() error {
	return d.err
}

func (v Values) copy() Values {
	values := make(Values, len(v))
	maps.Copy(values, v)

	return values
}
