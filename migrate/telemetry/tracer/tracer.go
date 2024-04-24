package tracer

import (
	"context"
	"net/url"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/alexfalkowski/migrieren/migrate/migrator"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// Migrator for otel.
type Migrator struct {
	migrator migrator.Migrator
	tracer   trace.Tracer
}

// NewMigrator for otel.
func NewMigrator(migrator migrator.Migrator, tracer trace.Tracer) *Migrator {
	return &Migrator{migrator: migrator, tracer: tracer}
}

// Migrate a database to a version and returning the database logs.
func (m *Migrator) Migrate(ctx context.Context, source, db string, version uint64) ([]string, error) {
	u, err := url.Parse(db)
	if err != nil {
		return nil, err
	}

	attrs := []attribute.KeyValue{
		semconv.DBSystemKey.String(u.Scheme),
		semconv.DBUser(u.User.Username()),
		attribute.Key("db.migrate.version").Int64(int64(version)),
	}

	ctx, span := m.tracer.Start(ctx, operationName("db"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))

	logs, err := m.migrator.Migrate(ctx, source, db, version)
	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return logs, err
}

// Ping the migrator.
func (m *Migrator) Ping(ctx context.Context, source, db string) error {
	return m.migrator.Ping(ctx, source, db)
}

func operationName(name string) string {
	return tracer.OperationName("migrate", name)
}
