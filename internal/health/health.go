package health

import (
	"github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/health"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/time"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/health/checker"
	"github.com/alexfalkowski/migrieren/internal/migrate"
)

// RegisterParams for health.
type RegisterParams struct {
	di.In
	Migrator *migrate.Migrator
	Server   *server.Server
	FS       *os.FS
	Migrate  *migrate.Config
	Config   *Config
	Name     env.Name
}

// Register for health.
func Register(params RegisterParams) {
	d := time.MustParseDuration(params.Config.Duration)
	t := time.MustParseDuration(params.Config.Timeout)
	regs := health.Registrations{
		server.NewRegistration("noop", d, checker.NewNoopChecker()),
		server.NewOnlineRegistration(d, d),
	}

	for _, db := range params.Migrate.Databases {
		checker := checker.NewMigrator(db, params.FS, params.Migrator, t)
		reg := server.NewRegistration(db.Name, d, checker)
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
