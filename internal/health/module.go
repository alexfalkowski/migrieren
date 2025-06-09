package health

import "go.uber.org/fx"

// Module for fx.
var Module = fx.Options(
	fx.Provide(registrations),
	fx.Provide(httpHealthObserver),
	fx.Provide(httpLivenessObserver),
	fx.Provide(httpReadinessObserver),
	fx.Provide(grpcObserver),
)
