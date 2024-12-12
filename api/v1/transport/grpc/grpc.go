package grpc

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/migrieren/api/migrate"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Register server.
func Register(gs *grpc.Server, server v1.ServiceServer) {
	v1.RegisterServiceServer(gs.Server(), server)
}

// NewServer for gRPC.
func NewServer(service *migrate.Migrator) v1.ServiceServer {
	return &Server{service: service}
}

// Server for gRPC.
type Server struct {
	v1.UnimplementedServiceServer
	service *migrate.Migrator
}

func (s *Server) error(err error) error {
	if migrate.IsNotFound(err) {
		return status.Error(codes.NotFound, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}
