package grpc

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/internal/api/v1/migrate"
	"github.com/alexfalkowski/migrieren/internal/diagnostics"
)

// Register registers server as the migrieren.v1 gRPC service implementation.
func Register(registrar grpc.ServiceRegistrar, server *Server) {
	v1.RegisterServiceServer(registrar, server)
}

// NewServer constructs a gRPC transport adapter around service.
func NewServer(service *migrate.Migrator) *Server {
	return &Server{migrator: service}
}

// Server implements the migrieren.v1 gRPC service.
//
// It delegates migration work to the versioned API migrator and maps domain
// errors to gRPC status codes.
type Server struct {
	v1.UnimplementedServiceServer
	migrator *migrate.Migrator
}

func (s *Server) error(err error) error {
	if migrate.IsInvalidVersion(err) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if migrate.IsNotFound(err) {
		return status.SafeError(codes.NotFound, err)
	}

	if migrate.IsCanceled(err) {
		return status.SafeError(codes.Canceled, err)
	}

	if migrate.IsDeadlineExceeded(err) {
		return status.SafeError(codes.DeadlineExceeded, err)
	}

	return status.SafeError(codes.Internal, err)
}

func setFailureTrailer(ctx context.Context, values diagnostics.Values) {
	_ = grpc.SetTrailer(ctx, meta.New(values.Map()))
}
