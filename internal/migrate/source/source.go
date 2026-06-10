package source

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/migrieren/internal/migrate/url"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

// Check validates that sourceURL can be used as a migration source without
// making unbounded external dependency calls.
func Check(ctx context.Context, sourceURL string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	u, err := url.Parse(sourceURL)
	if err != nil {
		return err
	}
	if u.Scheme == "github" {
		return nil
	}

	src, err := Open(sourceURL)
	if err != nil {
		return err
	}
	defer func() {
		_ = src.Close()
	}()

	return ctx.Err()
}

// Open opens a source driver.
func Open(sourceURL string) (source.Driver, error) {
	return source.Open(sourceURL)
}
