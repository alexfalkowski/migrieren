package grpc

import (
	"context"
	"fmt"

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
	config   *migrate.Config
	migrator migrator.Migrator
	v1.UnimplementedServiceServer
}

// Migrate for gRPC.
func (s *Server) Migrate(ctx context.Context, req *v1.MigrateRequest) (*v1.MigrateResponse, error) {
	resp := &v1.MigrateResponse{
		Migration: &v1.Migration{
			Database: req.Database,
			Version:  req.Version,
		},
	}

	d := s.config.Database(req.Database)
	if d == nil {
		return resp, status.Error(codes.NotFound, fmt.Sprintf("%s: not found", req.Database))
	}

	logs, err := s.migrator.Migrate(ctx, d.Source, d.URL, req.Version)
	if err != nil {
		return resp, status.Error(codes.Internal, err.Error())
	}

	resp.Meta = meta.Attributes(ctx)
	resp.Migration.Logs = logs

	return resp, nil
}
