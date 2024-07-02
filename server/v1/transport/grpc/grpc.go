package grpc

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/server/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Register server.
func Register(gs *grpc.Server, server v1.ServiceServer) {
	v1.RegisterServiceServer(gs.Server(), server)
}

// NewServer for gRPC.
func NewServer(service *service.Service) v1.ServiceServer {
	return &Server{service: service}
}

// Server for gRPC.
type Server struct {
	v1.UnimplementedServiceServer
	service *service.Service
}

func (s *Server) error(err error) error {
	if service.IsNotFound(err) {
		return status.Error(codes.NotFound, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}
