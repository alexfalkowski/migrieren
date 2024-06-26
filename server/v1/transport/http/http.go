package http

import (
	"github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/migrieren/server/service"
)

// Error for HTTP.
type Error struct {
	Message string `json:"message,omitempty"`
}

// Register for HTTP.
func Register(service *service.Service) {
	http.Handle("/v1/migrate", &migrateHandler{service: service})
}
