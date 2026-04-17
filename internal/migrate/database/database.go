package database

import (
	"errors"

	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/jackc/pgx/v5"
)

var (
	// ErrInvalidURL is returned when the url is invalid.
	ErrInvalidURL = errors.New("database: invalid url")

	// ErrUnsupportedDriver is returned when the driver is not supported.
	ErrUnsupportedDriver = errors.New("database: unsupported driver")
)

// Open opens a migrate database driver for databaseURL.
//
// URL handling errors are returned to the caller via the exported sentinel
// errors in this package.
//
// Telemetry wiring is treated differently on purpose: failures from
// telemetry.Open or telemetry.RegisterDBStatsMetrics are considered
// process-level misconfiguration/invariant violations for this service, so this
// function fails fast via runtime.Must rather than degrading to a runtime
// migration error.
func Open(databaseURL string) (database.Driver, error) {
	scheme, host, ok := splitURL(databaseURL)
	if !ok {
		return nil, ErrInvalidURL
	}

	switch scheme {
	case "pgx5":
		attrs := telemetry.WithAttributes(attributes.DBSystemNamePostgreSQL)

		db, err := telemetry.Open("pgx/v5", joinURL("postgres", host), attrs)
		// Fail fast: the service treats DB telemetry initialization as required.
		runtime.Must(err)

		_, err = telemetry.RegisterDBStatsMetrics(db, attrs)
		// Fail fast: running without DB stats metrics is an invalid process state.
		runtime.Must(err)

		return pgx.WithInstance(db, &pgx.Config{})
	default:
		return nil, ErrUnsupportedDriver
	}
}

func splitURL(url string) (string, string, bool) {
	return strings.Cut(url, "://")
}

func joinURL(scheme, host string) string {
	return strings.Join("://", scheme, host)
}
