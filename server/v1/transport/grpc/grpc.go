package grpc

import (
	"context"
	"fmt"

	sgrpc "github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/grpc/metrics/prometheus"
	"github.com/alexfalkowski/go-service/transport/grpc/trace/opentracing"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/version"
	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	"github.com/alexfalkowski/migrieren/migrate"
	"github.com/alexfalkowski/migrieren/migrate/migrator"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// RegisterParams for gRPC.
type RegisterParams struct {
	fx.In

	Lifecycle     fx.Lifecycle
	GRPCServer    *grpc.Server
	HTTPServer    *shttp.Server
	GRPCConfig    *sgrpc.Config
	Logger        *zap.Logger
	Tracer        opentracing.Tracer
	Version       version.Version
	Metrics       *prometheus.ClientMetrics
	MigrateConfig *migrate.Config
	Migrator      migrator.Migrator
}

// Register server.
func Register(params RegisterParams) {
	server := NewServer(ServerParams{Config: params.MigrateConfig, Migrator: params.Migrator})

	v1.RegisterServiceServer(params.GRPCServer, server)

	var conn *grpc.ClientConn

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			conn, _ = sgrpc.NewClient(
				sgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", params.GRPCConfig.Port), Version: params.Version, Config: params.GRPCConfig},
				sgrpc.WithClientLogger(params.Logger), sgrpc.WithClientTracer(params.Tracer), sgrpc.WithClientDialOption(grpc.WithBlock()),
				sgrpc.WithClientMetrics(params.Metrics),
			)

			return v1.RegisterServiceHandler(ctx, params.HTTPServer.Mux, conn)
		},
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})
}
