package grpc

import (
	"context"
	"fmt"

	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/http"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// RegisterParams for gRPC.
type RegisterParams struct {
	fx.In

	Lifecycle  fx.Lifecycle
	GRPCServer *grpc.Server
	HTTPServer *http.Server
	GRPCConfig *grpc.Config
	Logger     *zap.Logger
	Tracer     tracer.Tracer
	Meter      metric.Meter
	Server     v1.ServiceServer
}

// Register server.
func Register(params RegisterParams) error {
	ctx := context.Background()

	v1.RegisterServiceServer(params.GRPCServer.Server, params.Server)

	conn, err := grpc.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", params.GRPCConfig.Port), params.GRPCConfig,
		grpc.WithClientLogger(params.Logger), grpc.WithClientTracer(params.Tracer), grpc.WithClientMetrics(params.Meter),
		grpc.WithClientRetry(),
	)
	if err != nil {
		return err
	}

	if err := v1.RegisterServiceHandler(ctx, params.HTTPServer.Mux, conn); err != nil {
		return err
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			conn.ResetConnectBackoff()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})

	return nil
}
