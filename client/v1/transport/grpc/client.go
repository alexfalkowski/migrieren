package grpc

import (
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/security/token"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	v1c "github.com/alexfalkowski/migrieren/client/v1/config"
	"github.com/alexfalkowski/migrieren/transport/grpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ServiceClientParams for gRPC.
type ServiceClientParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Client    *v1c.Config
	Logger    *zap.Logger
	Tracer    trace.Tracer
	Meter     metric.Meter
	Generator token.Generator
	UserAgent env.UserAgent
}

// NewServiceClient for gRPC.
func NewServiceClient(params ServiceClientParams) (v1.ServiceClient, error) {
	opts := grpc.ClientOpts{
		Lifecycle: params.Lifecycle,
		Client:    params.Client.Config,
		Logger:    params.Logger,
		Tracer:    params.Tracer,
		Meter:     params.Meter,
		Generator: params.Generator,
		UserAgent: params.UserAgent,
	}
	conn, err := grpc.NewClient(opts)

	return v1.NewServiceClient(conn), err
}
