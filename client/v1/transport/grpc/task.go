package grpc

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/grpc/trace/opentracing"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/client/task"
	v1c "github.com/alexfalkowski/migrieren/client/v1/config"
	kzap "github.com/alexfalkowski/migrieren/client/v1/transport/grpc/logger/zap"
	gopentracing "github.com/alexfalkowski/migrieren/client/v1/transport/grpc/trace/opentracing"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// TaskParams for gRPC.
type TaskParams struct {
	fx.In

	Client v1.ServiceClient
	Config *v1c.Config
	Tracer opentracing.Tracer
	Logger *zap.Logger
}

// NewTask for gRPC.
func NewTask(params TaskParams) task.Task {
	var t task.Task = &Task{client: params.Client, cfg: params.Config}
	t = kzap.NewClient(params.Logger, params.Config, t)
	t = gopentracing.NewClient(params.Config, params.Tracer, t)

	return t
}

// Task for gRPC.
type Task struct {
	client v1.ServiceClient
	cfg    *v1c.Config
}

// Perform migrating the database.
func (t *Task) Perform(ctx context.Context) ([]string, error) {
	req := &v1.MigrateRequest{Database: t.cfg.Database, Version: t.cfg.Version}

	resp, err := t.client.Migrate(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Migration.Logs, nil
}
