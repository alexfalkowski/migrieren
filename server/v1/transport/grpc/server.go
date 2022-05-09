package grpc

import (
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
)

// ServerParams for gRPC.
type ServerParams struct {
}

// NewServer for gRPC.
func NewServer(params ServerParams) v1.ServiceServer {
	return &Server{}
}

// Server for gRPC.
type Server struct {
	v1.UnimplementedServiceServer
}
