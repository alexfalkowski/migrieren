package grpc

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/api/migrate"
)

// Register server.
func Register(registrar grpc.ServiceRegistrar, server *Server) {
	v1.RegisterServiceServer(registrar, server)
}

// NewServer for gRPC.
func NewServer(service *migrate.Migrator) *Server {
	return &Server{migrator: service}
}

// Server for gRPC.
type Server struct {
	v1.UnimplementedServiceServer
	migrator *migrate.Migrator
}

func (s *Server) error(err error) error {
	if err == nil {
		return nil
	}

	if migrate.IsNotFound(err) {
		return status.Error(codes.NotFound, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}
