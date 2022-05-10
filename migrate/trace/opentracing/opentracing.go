package opentracing

import (
	"context"
	"fmt"

	"github.com/alexfalkowski/go-service/trace/opentracing"
	"github.com/alexfalkowski/go-service/version"
	"github.com/alexfalkowski/migrieren/migrate/migrator"
	otr "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/fx"
)

// Tracer for opentracing.
type Tracer otr.Tracer

// StartSpanFromContext for opentracing.
func StartSpanFromContext(ctx context.Context, tracer Tracer, operation, method string, opts ...otr.StartSpanOption) (context.Context, otr.Span) {
	return opentracing.StartSpanFromContext(ctx, tracer, "migrator", operation, method, opts...)
}

// NewTracer for opentracing.
func NewTracer(lc fx.Lifecycle, cfg *opentracing.Config, version version.Version) (Tracer, error) {
	return opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Name: "migrator", Config: cfg, Version: version})
}

// Migrator for opentracing.
type Migrator struct {
	migrator migrator.Migrator
	tracer   Tracer
}

// NewMigrator for opentracing.
func NewMigrator(migrator migrator.Migrator, tracer Tracer) *Migrator {
	return &Migrator{migrator: migrator, tracer: tracer}
}

// Migrate a database to a version and returning the database logs.
func (m *Migrator) Migrate(ctx context.Context, source, db string, version uint64) ([]string, error) {
	ctx, span := StartSpanFromContext(ctx, m.tracer, "migrate", fmt.Sprintf("db to version %d", version))
	defer span.Finish()

	logs, err := m.migrator.Migrate(ctx, source, db, version)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.String("event", "error"), log.String("message", err.Error()))

		return nil, err
	}

	return logs, nil
}

// Ping the migrator.
func (m *Migrator) Ping(ctx context.Context, source, db string) error {
	return m.migrator.Ping(ctx, source, db)
}
