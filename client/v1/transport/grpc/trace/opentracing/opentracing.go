package opentracing

import (
	"context"
	"fmt"
	"strings"
	"time"

	stime "github.com/alexfalkowski/go-service/time"
	gopentracing "github.com/alexfalkowski/go-service/transport/grpc/trace/opentracing"
	"github.com/alexfalkowski/migrieren/client/task"
	v1 "github.com/alexfalkowski/migrieren/client/v1/config"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

// Client for opentracing.
type Client struct {
	cfg    *v1.Config
	tracer gopentracing.Tracer
	task.Task
}

// NewClient for zap.
func NewClient(cfg *v1.Config, tracer gopentracing.Tracer, task task.Task) *Client {
	return &Client{cfg: cfg, tracer: tracer, Task: task}
}

// Perform tracing for client.
func (c *Client) Perform(ctx context.Context) ([]string, error) {
	start := time.Now().UTC()
	operationName := fmt.Sprintf("sync %s/%d", c.cfg.Database, c.cfg.Version)
	opts := []opentracing.StartSpanOption{
		opentracing.Tag{Key: "client.start_time", Value: start.Format(time.RFC3339)},
		opentracing.Tag{Key: "client.database", Value: c.cfg.Database},
		opentracing.Tag{Key: "client.version", Value: c.cfg.Version},
		opentracing.Tag{Key: "component", Value: "client"},
		ext.SpanKindRPCClient,
	}

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, operationName, opts...)
	defer span.Finish()

	logs, err := c.Task.Perform(ctx)

	span.SetTag("client.duration", stime.ToMilliseconds(time.Since(start)))

	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.String("event", "error"), log.String("message", err.Error()))

		return nil, err
	}

	span.SetTag("client.logs", strings.Join(logs, ","))

	return logs, nil
}
