package grpc

import (
	"context"

	ac "github.com/alexfalkowski/auth/client"
	c "github.com/alexfalkowski/go-service/client"
	t "github.com/alexfalkowski/go-service/security/token"
	"github.com/alexfalkowski/go-service/transport/grpc"
	gt "github.com/alexfalkowski/go-service/transport/grpc/security/token"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	g "google.golang.org/grpc"
)

// ClientOpts for gRPC.
type ClientOpts struct {
	Lifecycle    fx.Lifecycle
	ClientConfig *c.Config
	TokenConfig  *t.Config
	Logger       *zap.Logger
	Tracer       trace.Tracer
	Meter        metric.Meter
	Token        *ac.Token
}

// NewClient for gRPC.
func NewClient(options ClientOpts) (*g.ClientConn, error) {
	sec, err := grpc.WithClientSecure(options.ClientConfig.Security)
	if err != nil {
		return nil, err
	}

	opts := []grpc.ClientOption{
		grpc.WithClientLogger(options.Logger), grpc.WithClientTracer(options.Tracer),
		grpc.WithClientMetrics(options.Meter), grpc.WithClientRetry(options.ClientConfig.Retry),
		grpc.WithClientUserAgent(options.ClientConfig.UserAgent), sec,
	}

	if IsAuth(options.TokenConfig) {
		opts = append(opts, grpc.WithClientDialOption(g.WithPerRPCCredentials(gt.NewPerRPCCredentials(options.Token.Generator("jwt", "migrieren")))))
	}

	conn, err := grpc.NewClient(options.ClientConfig.Host, opts...)

	options.Lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return conn.Close()
		},
	})

	return conn, err
}
