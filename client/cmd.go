package client

import (
	"context"

	"go.uber.org/fx"
)

// RunCommandParams for client.
type RunCommandParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Client    *Client
}

// RunCommand for client.
func RunCommand(params RunCommandParams) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := params.Client.Migrate(ctx)

			return err
		},
	})
}
