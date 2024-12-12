package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/rpc"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/migrieren/api/migrate"
)

// Register for HTTP.
func Register(service *migrate.Migrator) {
	mh := &migrateHandler{service: service}
	rpc.Route("/v1/migrate", mh.Migrate)
}

func handleError(err error) error {
	if migrate.IsNotFound(err) {
		return status.Error(http.StatusNotFound, err.Error())
	}

	return err
}
