package client

import (
	"context"

	v1 "github.com/alexfalkowski/migrieren/client/v1/config"
	"go.uber.org/fx"
)

// RunCommandParams for client.
type RunCommandParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *v1.Config
	Client    *Client
}

// RunCommand for client.
func RunCommand(params RunCommandParams) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := params.Client.Migrate(ctx, params.Config.Database, params.Config.Version)

			return err
		},
	})
}
