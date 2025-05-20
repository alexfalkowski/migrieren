package tracer

import (
	"context"
	"net/url"

	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/migrieren/internal/migrate/migrator"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// Migrator for otel.
type Migrator struct {
	migrator migrator.Migrator
	tracer   *tracer.Tracer
}

// NewMigrator for otel.
func NewMigrator(migrator migrator.Migrator, tracer *tracer.Tracer) *Migrator {
	return &Migrator{migrator: migrator, tracer: tracer}
}

// Migrate a database to a version and returning the database logs.
func (m *Migrator) Migrate(ctx context.Context, source, db string, version uint64) ([]string, error) {
	attrs := []attribute.KeyValue{
		semconv.DBSystemKey.String(m.system(db)),
		attribute.Key("db.migrate.version").Int64(int64(version)), //nolint:gosec
	}

	ctx, span := m.tracer.StartClient(ctx, operationName("db"), attrs...)
	defer span.End()

	logs, err := m.migrator.Migrate(ctx, source, db, version)
	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return logs, err
}

// Ping the migrator.
func (m *Migrator) Ping(ctx context.Context, source, db string) error {
	return m.migrator.Ping(ctx, source, db)
}

func (m *Migrator) system(db string) string {
	u, _ := url.Parse(db)
	if u != nil && !strings.IsEmpty(u.Scheme) {
		return u.Scheme
	}

	return db
}

func operationName(name string) string {
	return tracer.OperationName("migrate", name)
}
