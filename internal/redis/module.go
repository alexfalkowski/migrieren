package redis

import "github.com/alexfalkowski/go-service/v2/di"

// Module registers Redis dependencies for DI.
//
// It provides [NewClient] so packages depending on distributed locking can
// inject a [Client].
var Module = di.Module(
	di.Constructor(NewClient),
)
