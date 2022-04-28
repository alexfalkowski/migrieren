package health

import (
	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/health"
	mhealth "github.com/alexfalkowski/migrieren/health"
	"go.uber.org/fx"
)

// Params for health.
type Params struct {
	fx.In

	Health *mhealth.Config
}

// NewRegistrations for health.
func NewRegistrations(params Params) health.Registrations {
	registrations := health.Registrations{
		server.NewRegistration("noop", params.Health.Duration, checker.NewNoopChecker()),
	}

	return registrations
}
