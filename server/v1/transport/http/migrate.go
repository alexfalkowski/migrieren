package http

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
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
		Error     *Error            `json:"error,omitempty"`
		Migration *Migration        `json:"migration,omitempty"`
	}

	// Migration for a specific database and version with logs.
	Migration struct {
		Database string   `json:"database,omitempty"`
		Logs     []string `json:"logs,omitempty"`
		Version  uint64   `json:"version,omitempty"`
	}

	migrateErrorer struct{}
)

// Migrate for HTTP.
func (s *Server) Migrate(ctx context.Context, req *MigrateRequest) (*MigrateResponse, error) {
	resp := &MigrateResponse{
		Migration: &Migration{
			Database: req.Database,
			Version:  req.Version,
		},
	}

	logs, err := s.service.Migrate(ctx, req.Database, req.Version)
	if err != nil {
		return resp, err
	}

	resp.Migration.Logs = logs
	resp.Meta = meta.CamelStrings(ctx, "")

	return resp, nil
}

func (*migrateErrorer) Error(ctx context.Context, err error) *MigrateResponse {
	return &MigrateResponse{Meta: meta.CamelStrings(ctx, ""), Error: &Error{Message: err.Error()}}
}

func (*migrateErrorer) Status(err error) int {
	if service.IsNotFoundError(err) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
