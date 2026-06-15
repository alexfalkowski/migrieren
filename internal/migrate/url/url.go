package url

import (
	"fmt"
	"net/url"

	"github.com/alexfalkowski/go-service/v2/errors"
	migrate "github.com/golang-migrate/migrate/v4"
)

// ErrInvalid is returned when a migration URL cannot be parsed.
var ErrInvalid = errors.New("migrate url: invalid")

// URL is the parsed migration URL type.
type URL = url.URL

// Parse parses raw as a migration URL and rejects empty or schemeless values.
func Parse(raw string) (*URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	if u.Scheme == "" {
		return nil, ErrInvalid
	}

	return u, nil
}

// DatabaseURL converts a pgx5 migrate URL into the PostgreSQL URL consumed
// by the underlying SQL driver.
//
// It also removes golang-migrate custom query parameters, so migration-only
// options such as x-migrations-table are not passed to the SQL driver.
func DatabaseURL(u *URL) string {
	dbURL := *u
	dbURL.Scheme = "postgres"

	return migrate.FilterCustomQuery(&dbURL).String()
}
