package health

import "github.com/alexfalkowski/go-service/v2/di"

// Module for fx.
var Module = di.Module(
	di.Register(register),
	di.Register(httpHealthObserver),
	di.Register(httpLivenessObserver),
	di.Register(httpReadinessObserver),
	di.Register(grpcObserver),
)
