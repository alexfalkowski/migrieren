package health

import (
	health "github.com/alexfalkowski/go-health/server"
	grpc "github.com/alexfalkowski/go-service/v2/transport/grpc/health"
	http "github.com/alexfalkowski/go-service/v2/transport/http/health"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

func httpHealthObserver(server *health.Server, migrate *migrate.Config) *http.HealthObserver {
	names := []string{"online"}
	for _, d := range migrate.Databases {
		names = append(names, d.Name)
	}

	return &http.HealthObserver{Observer: server.Observe(names...)}
}

func httpLivenessObserver(server *health.Server) *http.LivenessObserver {
	return &http.LivenessObserver{Observer: server.Observe("noop")}
}

func httpReadinessObserver(server *health.Server) *http.ReadinessObserver {
	return &http.ReadinessObserver{Observer: server.Observe("noop")}
}

func grpcObserver(server *health.Server) *grpc.Observer {
	return &grpc.Observer{Observer: server.Observe("noop")}
}
