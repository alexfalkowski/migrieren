package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/rpc"
	"github.com/alexfalkowski/migrieren/server/migrate"
)

// Register for HTTP.
func Register(service *migrate.Migrator) {
	mh := &migrateHandler{service: service}
	rpc.Unary("/v1/migrate", mh.Migrate)
}

func handleError(err error) error {
	if migrate.IsNotFound(err) {
		return rpc.Error(http.StatusNotFound, err.Error())
	}

	return err
}
