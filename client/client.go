package client

import (
	"context"

	v1 "github.com/alexfalkowski/migrieren/api/migrieren/v1"
)

// Client for migrieren.
type Client struct {
	client v1.ServiceClient
}

// NewClient for migrieren.
func NewClient(client v1.ServiceClient) *Client {
	return &Client{client: client}
}

// Migrate the database to version.
func (c *Client) Migrate(ctx context.Context, database string, version uint64) ([]string, error) {
	req := &v1.MigrateRequest{Database: database, Version: version}

	resp, err := c.client.Migrate(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Migration.Logs, nil
}
