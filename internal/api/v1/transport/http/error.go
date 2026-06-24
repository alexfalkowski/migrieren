package http

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/migrieren/internal/api/v1/migrate"
)

func responseError(err error) error {
	if migrate.IsInvalidVersion(err) {
		return status.Error(http.StatusBadRequest, err.Error())
	}

	if migrate.IsNotFound(err) {
		return status.SafeError(http.StatusNotFound, err)
	}

	if migrate.IsCanceled(err) {
		return status.SafeError(http.StatusClientClosedRequest, err)
	}

	if migrate.IsDeadlineExceeded(err) {
		return status.SafeError(http.StatusGatewayTimeout, err)
	}

	return status.SafeError(http.StatusInternalServerError, err)
}
