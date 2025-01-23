package http

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
)

type (
	// MigrateRequest for a specific database and version.
	MigrateRequest struct {
		Database string `json:"database,omitempty"`
		Version  uint64 `json:"version,omitempty"`
	}

	// MigrateResponse for a specific database and version.
	MigrateResponse struct {
		Meta      map[string]string `json:"meta,omitempty"`
		Migration *Migration        `json:"migration,omitempty"`
	}

	// Migration for a specific database and version with logs.
	Migration struct {
		Database string   `json:"database,omitempty"`
		Logs     []string `json:"logs,omitempty"`
		Version  uint64   `json:"version,omitempty"`
	}
)

// Migrate for HTTP.
func (h *Handler) Migrate(ctx context.Context, req *MigrateRequest) (*MigrateResponse, error) {
	resp := &MigrateResponse{
		Migration: &Migration{
			Database: req.Database,
			Version:  req.Version,
		},
	}
	logs, err := h.migrator.Migrate(ctx, req.Database, req.Version)

	resp.Migration.Logs = logs
	resp.Meta = meta.CamelStrings(ctx, "")

	return resp, h.error(err)
}
