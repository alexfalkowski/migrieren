package migrate

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// IsNotFound reports whether err indicates that the requested database name was
// not present in the configured database list.
func IsNotFound(err error) bool {
	return errors.Is(err, migrate.ErrNotFound)
}

// IsCanceled reports whether err indicates the caller canceled migration work.
func IsCanceled(err error) bool {
	return errors.Is(err, migrate.ErrMigrationCanceled)
}

// IsDeadlineExceeded reports whether err indicates migration work exceeded the
// request deadline.
func IsDeadlineExceeded(err error) bool {
	return errors.Is(err, migrate.ErrMigrationDeadlineExceeded)
}
