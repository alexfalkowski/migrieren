package http

import (
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/health/transport/http"
	"github.com/alexfalkowski/migrieren/migrate"
)

// NewHealthObserver for HTTP.
func NewHealthObserver(healthServer *server.Server, migrate *migrate.Config) *http.HealthObserver {
	names := make([]string, len(migrate.Databases))
	for i, d := range migrate.Databases {
		names[i] = d.Name
	}

	return &http.HealthObserver{Observer: healthServer.Observe(names...)}
}

// NewLivenessObserver for HTTP.
func NewLivenessObserver(healthServer *server.Server) *http.LivenessObserver {
	return &http.LivenessObserver{Observer: healthServer.Observe("noop")}
}

// NewReadinessObserver for HTTP.
func NewReadinessObserver(healthServer *server.Server) *http.ReadinessObserver {
	return &http.ReadinessObserver{Observer: healthServer.Observe("noop")}
}
