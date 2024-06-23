package http

import (
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/migrieren/server/service"
	"go.uber.org/fx"
)

type (
	// RegisterParams for HTTP.
	RegisterParams struct {
		fx.In

		Marshaller *marshaller.Map
		Mux        http.ServeMux
		Service    *service.Service
	}

	// Server for HTTP.
	Server struct {
		service *service.Service
	}

	// Error for HTTP.
	Error struct {
		Message string `json:"message,omitempty"`
	}
)

// Register for HTTP.
func Register(params RegisterParams) error {
	s := &Server{service: params.Service}

	mh := http.NewHandler[MigrateRequest](params.Mux, params.Marshaller, &migrateErrorer{})
	if err := mh.Handle("POST", "/v1/migrate", s.Migrate); err != nil {
		return err
	}

	return nil
}
