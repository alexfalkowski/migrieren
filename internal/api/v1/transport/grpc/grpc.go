package grpc

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/api/migrate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Register server.
func Register(gs *grpc.Server, server *Server) {
	v1.RegisterServiceServer(gs.Server(), server)
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
