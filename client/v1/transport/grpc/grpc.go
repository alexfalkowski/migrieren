package grpc

import (
	"context"

	"github.com/alexfalkowski/migrieren/client/task"
	"go.uber.org/fx"
)

// RegisterParams for gRPC.
type RegisterParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Task      task.Task
}

// Register client.
func Register(params RegisterParams) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := params.Task.Perform(ctx)

			return err
		},
	})
}
