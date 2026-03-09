package health

import "github.com/alexfalkowski/go-service/v2/di"

// Module registers health checks and endpoint observers for DI.
var Module = di.Module(
	di.Register(Register),
	di.Register(httpHealthObserver),
	di.Register(httpLivenessObserver),
	di.Register(httpReadinessObserver),
	di.Register(grpcObserver),
)
