package client

import (
	v1 "github.com/alexfalkowski/migrieren/client/v1"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	v1.Module,
	fx.Provide(NewClient),
)
