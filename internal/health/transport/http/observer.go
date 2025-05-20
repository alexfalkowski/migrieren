package http

import (
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/v2/health/transport/http"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// NewHealthObserver for HTTP.
func NewHealthObserver(server *server.Server, migrate *migrate.Config) *http.HealthObserver {
	names := []string{"online"}
	for _, d := range migrate.Databases {
		names = append(names, d.Name)
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
