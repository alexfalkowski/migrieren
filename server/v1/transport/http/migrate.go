package http

import (
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/migrieren/server/service"
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

	migrateHandler struct {
		service *service.Service
	}
)

func (h *migrateHandler) Handle(ctx http.Context, req *MigrateRequest) (*MigrateResponse, error) {
	resp := &MigrateResponse{
		Migration: &Migration{
			Database: req.Database,
			Version:  req.Version,
		},
	}

	logs, err := h.service.Migrate(ctx, req.Database, req.Version)
	if err != nil {
		resp.Meta = meta.CamelStrings(ctx, "")

		return resp, handleError(err)
	}

	resp.Migration.Logs = logs
	resp.Meta = meta.CamelStrings(ctx, "")

	return resp, nil
}
