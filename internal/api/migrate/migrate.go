package migrate

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

const (
	// FailureStageKey is the safe metadata key used to identify setup-stage
	// migration failures for transport diagnostics.
	FailureStageKey = "migrateErrorStage"

	// FailureStageSource identifies failures while resolving a migration source.
	FailureStageSource = "source"

	// FailureStageURL identifies failures while resolving a database URL.
	FailureStageURL = "url"
)

// IsNotFound reports whether err indicates that the requested database name was
// not present in the configured database list.
//
// This is intended for transport layers to map the condition to an appropriate
// status code (for example gRPC NotFound / HTTP 404).
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

// IsInvalidConfig reports whether err indicates invalid migration configuration.
func IsInvalidConfig(err error) bool {
	return errors.Is(err, migrate.ErrInvalidConfig)
}

// IsInvalidMigration reports whether err indicates a migration execution failure.
func IsInvalidMigration(err error) bool {
	return errors.Is(err, migrate.ErrInvalidMigration)
}

// FailureKind returns a stable, safe diagnostic kind for err.
func FailureKind(ctx context.Context, err error) string {
	switch {
	case IsNotFound(err):
		return "not_found"
	case IsCanceled(err):
		return "canceled"
	case IsDeadlineExceeded(err):
		return "deadline_exceeded"
	case !meta.Attribute(ctx, FailureStageKey).IsEmpty(), IsInvalidConfig(err):
		return "invalid_config"
	case IsInvalidMigration(err):
		return "invalid_migration"
	default:
		return "unknown"
	}
}

// NewMigrator constructs a transport-facing [Migrator].
//
// Dependencies:
//   - migrator: the core migrator that executes migrations given a source URL and
//     database URL.
//   - fs: a filesystem abstraction used to resolve `Database.Source` and
//     `Database.URL` values (for example `file:...`).
//   - cfg: the migration configuration containing the list of named databases.
func NewMigrator(migrator *migrate.Migrator, fs *os.FS, cfg *migrate.Config) *Migrator {
	return &Migrator{migrator: migrator, fs: fs, config: cfg}
}

// Migrator adapts the core migrator to a "database name + version" API that is
// convenient for transport layers.
//
// The adapter:
//   - looks up a database by name in the provided config,
//   - reads its source and URL through the filesystem abstraction,
//   - delegates the actual migration execution to the core migrator.
type Migrator struct {
	migrator *migrate.Migrator
	config   *migrate.Config
	fs       *os.FS
}

// Migrate migrates the named database to the given target version.
//
// The database is resolved from configuration by name. Its migration source and
// database URL are read via the filesystem abstraction, then passed to the core
// migrator.
//
// Returns the input context, or a derived context when the core migrator adds
// metadata, plus migration logs from the core migrator. If the database name
// does not exist in the configuration, this returns an error that wraps
// `internal/migrate.ErrNotFound` (detectable via [IsNotFound]).
func (s *Migrator) Migrate(ctx context.Context, db string, version uint64) (context.Context, []string, error) {
	d, err := s.config.Database(db)
	if d == nil {
		return ctx, nil, fmt.Errorf("%s: %w", db, err)
	}

	source, err := d.GetSource(s.fs)
	if err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair(FailureStageKey, meta.String(FailureStageSource)))
		return ctx, nil, err
	}

	url, err := d.GetURL(s.fs)
	if err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair(FailureStageKey, meta.String(FailureStageURL)))
		return ctx, nil, err
	}

	return s.migrator.Migrate(ctx, bytes.String(source), bytes.String(url), version)
}
