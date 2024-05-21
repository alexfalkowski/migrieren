package migrate

import (
	"context"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/migrieren/client"
	"go.uber.org/fx"
)

// Start for client.
func Start(lc fx.Lifecycle, client *client.Client) {
	cmd.Start(lc, func(ctx context.Context) {
		_, err := client.Migrate(ctx)
		runtime.Must(err)
	})
}
