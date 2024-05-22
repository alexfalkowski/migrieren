package client

import (
	"context"

	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
	v1c "github.com/alexfalkowski/migrieren/client/v1/config"
)

// Client for migrieren.
type Client struct {
	client v1.ServiceClient
	config *v1c.Config
}

// NewClient for migrieren.
func NewClient(client v1.ServiceClient, config *v1c.Config) *Client {
	return &Client{client: client, config: config}
}

// Migrate the database to version.
func (c *Client) Migrate(ctx context.Context) ([]string, error) {
	cfg := c.config.Migrate
	req := &v1.MigrateRequest{Database: cfg.Database, Version: cfg.Version}

	resp, err := c.client.Migrate(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetMigration().GetLogs(), nil
}
