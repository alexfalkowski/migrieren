package health

import (
	"github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/health"
	"github.com/alexfalkowski/go-service/v2/os"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/health/checker"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// RegisterParams contains the dependencies required to register health checks
// and observers.
type RegisterParams struct {
	di.In
	Migrator *migrate.Migrator
	Server   *server.Server
	FS       *os.FS
	Migrate  *migrate.Config
	Config   *Config
	Name     env.Name
}

// Register installs service-level and per-database health registrations.
//
// The application service name receives the noop, online, and per-database
// checks. The gRPC service name is registered with the noop check only so gRPC
// health remains healthy even when intentionally invalid database entries exist
// in test configuration.
func Register(params RegisterParams) {
	regs := health.Registrations{
		server.NewRegistration("noop", params.Config.Duration.Duration(), checker.NewNoopChecker()),
		server.NewOnlineRegistration(params.Config.Timeout.Duration(), params.Config.Duration.Duration()),
	}

	for _, db := range params.Migrate.Databases {
		checker := checker.NewMigrator(db, params.FS, params.Migrator, params.Config.Timeout)
		reg := server.NewRegistration(db.Name, params.Config.Duration.Duration(), checker)
		regs = append(regs, reg)
	}

	params.Server.Register(params.Name.String(), regs...)
	params.Server.Register(v1.Service_ServiceDesc.ServiceName, regs[0])
}

func httpHealthObserver(name env.Name, server *server.Server, migrate *migrate.Config) error {
	names := make([]string, 0, len(migrate.Databases)+1)
	names = append(names, "online")
	for _, d := range migrate.Databases {
		names = append(names, d.Name)
	}

	return server.Observe(name.String(), "healthz", names...)
}

func httpLivenessObserver(name env.Name, server *server.Server) error {
	return server.Observe(name.String(), "livez", "noop")
}

func httpReadinessObserver(name env.Name, server *server.Server) error {
	return server.Observe(name.String(), "readyz", "noop")
}

func grpcObserver(server *server.Server) error {
	return server.Observe(v1.Service_ServiceDesc.ServiceName, "grpc", "noop")
}
