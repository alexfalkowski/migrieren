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

// RegisterParams contains dependencies required to register health checks.
type RegisterParams struct {
	di.In
	// Migrator performs migration source/database validation for DB checks.
	Migrator *migrate.Migrator
	// Server is the go-health server used to register check sets.
	Server *server.Server
	// FS resolves source/url references in migration database config.
	FS *os.FS
	// Migrate contains configured databases that should receive health checks.
	Migrate *migrate.Config
	// Config contains interval and timeout values used by the scheduler/checkers.
	Config *Config
	// Name is the service name used for observer registration.
	Name env.Name
}

// Register registers health checks for the service.
//
// It parses health durations from params.Config, registers:
//   - a noop check,
//   - an online check,
//   - one migrator checker per configured database.
//
// Checks are registered under the process service name. In addition, the noop
// check is registered under the gRPC service name for transport-specific
// liveness integration.
//
// Panics:
//   - invalid Config.Duration or Config.Timeout values cause panic through
//     time.MustParseDuration.
func Register(params RegisterParams) {
	d := time.MustParseDuration(params.Config.Duration)
	t := time.MustParseDuration(params.Config.Timeout)
	regs := health.Registrations{
		server.NewRegistration("noop", d, checker.NewNoopChecker()),
		server.NewOnlineRegistration(t, d),
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
