package client

import (
	"context"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/runtime"
	"go.uber.org/fx"
)

// Migrate for client.
func Migrate(lc fx.Lifecycle, client *Client) {
	cmd.Start(lc, func(ctx context.Context) {
		_, err := client.Migrate(ctx)
		runtime.Must(err)
	})
}
