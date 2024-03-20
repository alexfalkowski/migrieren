package tracer

import (
	"context"
	"net/url"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/alexfalkowski/go-service/version"
	"github.com/alexfalkowski/migrieren/migrate/migrator"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// Tracer for otel.
type Tracer trace.Tracer

// NewTracer for otel.
func NewTracer(lc fx.Lifecycle, cfg *tracer.Config, env env.Environment, ver version.Version) (Tracer, error) {
	return tracer.NewTracer(context.Background(), lc, "migrator", env, ver, cfg)
}

// Migrator for otel.
type Migrator struct {
	migrator migrator.Migrator
	tracer   Tracer
}

// NewMigrator for otel.
func NewMigrator(migrator migrator.Migrator, tracer Tracer) *Migrator {
	return &Migrator{migrator: migrator, tracer: tracer}
}

// Migrate a database to a version and returning the database logs.
func (m *Migrator) Migrate(ctx context.Context, source, db string, version uint64) ([]string, error) {
	u, err := url.Parse(db)
	if err != nil {
		return nil, err
	}

	operationName := "migrate"
	attrs := []attribute.KeyValue{
		semconv.DBSystemKey.String(u.Scheme),
		semconv.DBUser(u.User.Username()),
		attribute.Key("db.migrate.version").Int64(int64(version)),
	}

	ctx, span := m.tracer.Start(
		ctx,
		operationName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))

	logs, err := m.migrator.Migrate(ctx, source, db, version)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return nil, err
	}

	return logs, nil
}

// Ping the migrator.
func (m *Migrator) Ping(ctx context.Context, source, db string) error {
	return m.migrator.Ping(ctx, source, db)
}
