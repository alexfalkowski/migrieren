package health

import "github.com/alexfalkowski/go-service/v2/di"

// Module for fx.
var Module = di.Module(
	di.Constructor(registrations),
	di.Constructor(httpHealthObserver),
	di.Constructor(httpLivenessObserver),
	di.Constructor(httpReadinessObserver),
	di.Constructor(grpcObserver),
)
