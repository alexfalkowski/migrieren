package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/rpc"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/migrieren/api/migrate"
)

// Register for HTTP.
func Register(handler *Handler) {
	rpc.Route("/v1/migrate", handler.Migrate)
}

// NewHandler for HTTP.
func NewHandler(migrator *migrate.Migrator) *Handler {
	return &Handler{migrator: migrator}
}

// Handler for HTTP.
type Handler struct {
	migrator *migrate.Migrator
}

func (h *Handler) error(err error) error {
	if err == nil {
		return nil
	}

	if migrate.IsNotFound(err) {
		return status.Error(http.StatusNotFound, err.Error())
	}

	return err
}
