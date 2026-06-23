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
//
// The returned value is one of "not_found", "canceled", "deadline_exceeded",
// "invalid_config", "invalid_migration", or "unknown". A context carrying
// [FailureStageKey] is classified as "invalid_config" so source/URL resolution
// failures are reported consistently by transports.
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
//   - delegates migration execution or status inspection to the core migrator,
//   - lists configured logical database names without exposing source or URL values.
type Migrator struct {
	migrator *migrate.Migrator
	config   *migrate.Config
	fs       *os.FS
}

// Databases returns the configured logical database names in config order.
//
// It intentionally exposes only names, not configured source strings, database
// URL strings, or resolved secret values.
func (s *Migrator) Databases() []string {
	databases := make([]string, 0, len(s.config.Databases))
	for _, d := range s.config.Databases {
		databases = append(databases, d.Name)
	}

	return databases
}

// Migrate migrates the named database to the given target version.
//
// The database is resolved from configuration by name. Its migration source and
// database URL are read via the filesystem abstraction, then passed to the core
// migrator.
//
// This adapter does not perform public request validation such as rejecting a
// zero version; transport callers that expose the migrieren.v1 API must enforce
// that contract before calling Migrate.
//
// Returns the input context, or a derived context when the core migrator adds
// metadata, plus migration logs from the core migrator. If the database name
// does not exist in the configuration, this returns an error that wraps
// `internal/migrate.ErrNotFound` (detectable via [IsNotFound]).
//
// Source and URL resolution failures return the underlying filesystem error and
// a derived context containing [FailureStageKey] set to [FailureStageSource] or
// [FailureStageURL]. [FailureKind] maps those staged failures to
// "invalid_config" for transport diagnostics.
func (s *Migrator) Migrate(ctx context.Context, db string, version uint64) (context.Context, []string, error) {
	d, err := s.config.Database(db)
	if err != nil {
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

// Status reports the current migration version state for the named database.
//
// The database is resolved from configuration by name. Only its database URL is
// read via the filesystem abstraction; source resolution is not needed for this
// status inspection path.
//
// If the database name does not exist in the configuration, this returns an
// error that wraps `internal/migrate.ErrNotFound` (detectable via [IsNotFound]).
// URL resolution failures return the underlying filesystem error and a derived
// context containing [FailureStageKey] set to [FailureStageURL].
func (s *Migrator) Status(ctx context.Context, db string) (context.Context, *migrate.Status, error) {
	d, err := s.config.Database(db)
	if err != nil {
		return ctx, nil, fmt.Errorf("%s: %w", db, err)
	}

	url, err := d.GetURL(s.fs)
	if err != nil {
		ctx = meta.WithAttributes(ctx, meta.NewPair(FailureStageKey, meta.String(FailureStageURL)))
		return ctx, nil, err
	}

	return s.migrator.Status(ctx, bytes.String(url))
}
