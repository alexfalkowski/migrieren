package http

import (
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/health/transport/http"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// NewHealthObserver for HTTP.
func NewHealthObserver(server *server.Server, migrate *migrate.Config) *http.HealthObserver {
	names := make([]string, len(migrate.Databases))
	for i, d := range migrate.Databases {
		names[i] = d.Name
	}

	return &http.HealthObserver{Observer: server.Observe(names...)}
}

// NewLivenessObserver for HTTP.
func NewLivenessObserver(server *server.Server) *http.LivenessObserver {
	return &http.LivenessObserver{Observer: server.Observe("noop")}
}

// NewReadinessObserver for HTTP.
func NewReadinessObserver(server *server.Server) *http.ReadinessObserver {
	return &http.ReadinessObserver{Observer: server.Observe("noop")}
}
