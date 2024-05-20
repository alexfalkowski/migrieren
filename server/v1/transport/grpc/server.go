package grpc

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/migrate"
	"github.com/alexfalkowski/migrieren/migrate/migrator"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ServerParams for gRPC.
type ServerParams struct {
	fx.In

	Config   *migrate.Config
	Migrator migrator.Migrator
}

// NewServer for gRPC.
func NewServer(params ServerParams) v1.ServiceServer {
	return &Server{config: params.Config, migrator: params.Migrator}
}

// Server for gRPC.
type Server struct {
	v1.UnimplementedServiceServer
	migrator migrator.Migrator
	config   *migrate.Config
}

// Migrate for gRPC.
func (s *Server) Migrate(ctx context.Context, req *v1.MigrateRequest) (*v1.MigrateResponse, error) {
	db := req.GetDatabase()
	ver := req.GetVersion()
	resp := &v1.MigrateResponse{
		Migration: &v1.Migration{
			Database: db,
			Version:  ver,
		},
	}

	d := s.config.Database(db)
	if d == nil {
		return resp, status.Error(codes.NotFound, db+": not found")
	}

	u, err := d.GetURL()
	if err != nil {
		return resp, status.Error(codes.Internal, err.Error())
	}

	logs, err := s.migrator.Migrate(ctx, d.Source, u, ver)
	if err != nil {
		return resp, status.Error(codes.Internal, err.Error())
	}

	resp.Meta = s.meta(ctx)
	resp.Migration.Logs = logs

	return resp, nil
}

func (s *Server) meta(ctx context.Context) map[string]string {
	return meta.CamelStrings(ctx, "")
}
