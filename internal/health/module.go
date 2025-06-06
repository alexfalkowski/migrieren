package health

import (
	"github.com/alexfalkowski/migrieren/internal/health/transport/grpc"
	"github.com/alexfalkowski/migrieren/internal/health/transport/http"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(http.NewHealthObserver),
	fx.Provide(http.NewLivenessObserver),
	fx.Provide(http.NewReadinessObserver),
	fx.Provide(grpc.NewObserver),
	fx.Provide(NewRegistrations),
)
