package diagnostics

import (
	"maps"
	"strconv"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

const invalidConfig = "invalid_config"

// StageSource identifies failures while resolving a migration source.
const StageSource = "source"

// StageURL identifies failures while resolving a database URL.
const StageURL = "url"

// Error wraps err with safe migration diagnostic values.
//
// The returned error unwraps to err, so errors.Is and errors.As continue to
// match the original cause.
func Error(err error, logs []string) error {
	return &diagnosticError{
		err:    err,
		values: newValues(err, logs),
	}
}

// InvalidConfig wraps err with diagnostics for a migration configuration
// resolution failure.
func InvalidConfig(err error, stage string) error {
	values := newValues(err, nil)
	values["migration-error"] = invalidConfig
	if stage != "" {
		values["migration-stage"] = stage
	}

	return &diagnosticError{
		err:    err,
		values: values,
	}
}

func newValues(err error, logs []string) Values {
	values := Values{
		"migration-error":     failureKind(err),
		"migration-log-count": strconv.Itoa(len(logs)),
	}

	if len(logs) > 0 {
		values["migration-log-last"] = logs[len(logs)-1]
	}

	return values
}

func failureKind(err error) string {
	switch {
	case errors.Is(err, migrate.ErrNotFound):
		return "not_found"
	case errors.Is(err, migrate.ErrMigrationCanceled):
		return "canceled"
	case errors.Is(err, migrate.ErrMigrationDeadlineExceeded):
		return "deadline_exceeded"
	case errors.Is(err, migrate.ErrInvalidConfig):
		return invalidConfig
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

// Values contains safe diagnostic key/value pairs attached to an error.
type Values map[string]string

// Map returns a copy of v as a plain string map.
func (v Values) Map() map[string]string {
	values := make(map[string]string, len(v))
	maps.Copy(values, v)

	return values
}

func (v Values) copy() Values {
	return Values(v.Map())
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
