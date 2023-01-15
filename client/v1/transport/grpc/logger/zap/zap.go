package zap

import (
	"context"
	"strings"
	"time"

	stime "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/migrieren/client/task"
	v1 "github.com/alexfalkowski/migrieren/client/v1/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Client for zap.
type Client struct {
	logger *zap.Logger
	cfg    *v1.Config
	task.Task
}

// NewClient for zap.
func NewClient(logger *zap.Logger, cfg *v1.Config, task task.Task) *Client {
	return &Client{logger: logger, cfg: cfg, Task: task}
}

// Perform logger for client.
func (c *Client) Perform(ctx context.Context) ([]string, error) {
	start := time.Now().UTC()
	logs, err := c.Task.Perform(ctx)
	fields := []zapcore.Field{
		zap.Int64("client.duration", stime.ToMilliseconds(time.Since(start))),
		zap.String("client.start_time", start.Format(time.RFC3339)),
		zap.String("client.database", c.cfg.Database),
		zap.Uint64("client.version", c.cfg.Version),
		zap.String("span.kind", "client"),
		zap.String("component", "client"),
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String("client.deadline", d.UTC().Format(time.RFC3339)))
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		c.logger.Error("finished call with error", fields...)

		return nil, err
	}

	fields = append(fields, zap.String("client.logs", strings.Join(logs, ",")))

	c.logger.Info("finished call with success", fields...)

	return logs, nil
}
