package health

import (
	"context"

	"github.com/alexfalkowski/go-health/v2/checker"
	"github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/health"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/time"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/migrate"
	"github.com/alexfalkowski/migrieren/internal/migrate/migrator"
)

func register(name env.Name, srv *server.Server, mig migrator.Migrator, fs *os.FS, migrate *migrate.Config, cfg *Config) {
	d := time.MustParseDuration(cfg.Duration)
	regs := health.Registrations{
		server.NewRegistration("noop", d, checker.NewNoopChecker()),
		server.NewOnlineRegistration(d, d),
	}

	for _, db := range migrate.Databases {
		checker := &migratorChecker{db: db, fs: fs, migrator: mig}
		reg := server.NewRegistration(db.Name, d, checker)
		regs = append(regs, reg)
	}

	srv.Register(name.String(), regs...)
	srv.Register(v1.Service_ServiceDesc.ServiceName, regs...)
}

func httpHealthObserver(name env.Name, server *server.Server, migrate *migrate.Config) error {
	names := []string{"online"}
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

type migratorChecker struct {
	db       *migrate.Database
	fs       *os.FS
	migrator migrator.Migrator
}

// Check the migrator.
func (c *migratorChecker) Check(ctx context.Context) error {
	source, err := c.db.GetSource(c.fs)
	if err != nil {
		return err
	}

	url, err := c.db.GetURL(c.fs)
	if err != nil {
		return err
	}

	return c.migrator.Ping(ctx, bytes.String(source), bytes.String(url))
}
