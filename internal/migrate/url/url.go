package url

import (
	"fmt"
	"net/url"

	"github.com/alexfalkowski/go-service/v2/errors"
	migrate "github.com/golang-migrate/migrate/v4"
)

// ErrInvalid is returned when a migration URL cannot be parsed.
var ErrInvalid = errors.New("migrate url: invalid")

// Parse parses raw as a migration URL and rejects empty or schemeless values.
func Parse(raw string) (*url.URL, error) {
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
func DatabaseURL(u *url.URL) string {
	dbURL := *u
	dbURL.Scheme = "postgres"

	return migrate.FilterCustomQuery(&dbURL).String()
}
