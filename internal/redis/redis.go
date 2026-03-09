package redis

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	otel "github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

// Client is an alias for redsync.Redsync.
//
// It is used to create distributed mutexes around migration operations.
type Client = redsync.Redsync

// NewClient constructs a distributed lock client backed by Redis.
//
// It reads Redis URL data from cfg.URL using fs.ReadSource, parses the URL into
// go-redis options, disables maintenance notifications, and enables Redis
// OpenTelemetry tracing and metrics instrumentation.
//
// Returns a [Client] on success.
//
// Errors:
//   - returns an underlying source read error if cfg.URL cannot be resolved.
//   - returns an underlying parse error if URL data is not a valid Redis URL.
//
// Panics:
//   - instrumentation setup is wrapped in runtime.Must and panics if tracing or
//     metrics hooks cannot be registered.
func NewClient(fs *os.FS, cfg *Config) (*Client, error) {
	data, err := fs.ReadSource(cfg.URL)
	if err != nil {
		return nil, err
	}

	opts, err := redis.ParseURL(bytes.String(data))
	if err != nil {
		return nil, err
	}

	opts.MaintNotificationsConfig = &maintnotifications.Config{
		Mode: maintnotifications.ModeDisabled,
	}

	client := redis.NewClient(opts)
	runtime.Must(otel.InstrumentTracing(client))
	runtime.Must(otel.InstrumentMetrics(client))

	return redsync.New(goredis.NewPool(client)), nil
}
