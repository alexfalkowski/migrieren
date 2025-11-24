package database

import (
	"errors"

	"github.com/XSAM/otelsql"
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	pgxstdlib "github.com/jackc/pgx/v5/stdlib"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

var (
	// ErrInvalidURL is returned when the url is invalid.
	ErrInvalidURL = errors.New("database: invalid url")

	// ErrUnsupportedDriver is returned when the driver is not supported.
	ErrUnsupportedDriver = errors.New("database: unsupported driver")
)

// Register registers the database that are supported.
func Register() {
	_ = driver.Register("postgres", pgxstdlib.GetDefaultDriver())
}

// Open from the specified URL.
func Open(databaseURL string) (database.Driver, error) {
	scheme, host, ok := splitURL(databaseURL)
	if !ok {
		return nil, ErrInvalidURL
	}

	switch scheme {
	case "pgx5":
		db, err := otelsql.Open("postgres", joinURL("postgres", host), otelsql.WithAttributes(semconv.DBSystemNamePostgreSQL))
		runtime.Must(err)

		err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(semconv.DBSystemNamePostgreSQL))
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
