package health

import "github.com/alexfalkowski/go-signal/v2"

// Module wires health registrations and observers into the DI container.
var Module = signal.Module(
	signal.Register(Register),
	signal.Register(httpHealthObserver),
	signal.Register(httpLivenessObserver),
	signal.Register(httpReadinessObserver),
	signal.Register(grpcObserver),
)
