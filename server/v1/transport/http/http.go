package http

import (
	"github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/migrieren/server/service"
)

// Register for HTTP.
func Register(service *service.Service) {
	http.Handle("/v1/migrate", &migrateHandler{service: service})
}
