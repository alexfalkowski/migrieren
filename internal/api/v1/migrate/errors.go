package migrate

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/migrieren/internal/api/migrate"
)

// ErrInvalidVersion is returned when a v1 migrate request uses a version outside
// the public API range.
var ErrInvalidVersion = errors.New("version must be between 1 and math.MaxInt")

// IsInvalidVersion reports whether err indicates an unsupported target version.
func IsInvalidVersion(err error) bool {
	return errors.Is(err, ErrInvalidVersion)
}

// IsNotFound reports whether err indicates that the requested database name was
// not present in the configured database list.
func IsNotFound(err error) bool {
	return migrate.IsNotFound(err)
}

// IsCanceled reports whether err indicates the caller canceled migration work.
func IsCanceled(err error) bool {
	return migrate.IsCanceled(err)
}

// IsDeadlineExceeded reports whether err indicates migration work exceeded the
// request deadline.
func IsDeadlineExceeded(err error) bool {
	return migrate.IsDeadlineExceeded(err)
}
