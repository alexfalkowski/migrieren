package http

import (
	"net/http"

	nh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/net/http/rpc"
	"github.com/alexfalkowski/migrieren/server/migrate"
)

// Register for HTTP.
func Register(service *migrate.Migrator) {
	rpc.Handle("/v1/migrate", &migrateHandler{service: service})
}

func handleError(err error) error {
	if migrate.IsNotFound(err) {
		return nh.Error(http.StatusNotFound, err.Error())
	}

	return err
}
