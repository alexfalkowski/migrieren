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
//
// GitHub sources are parsed but not opened here, because the upstream source
// driver does not expose a request-scoped timeout hook. The actual GitHub source
// is opened later during migration execution.
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

// Open opens a source driver using the upstream golang-migrate source registry.
//
// This package registers the file and GitHub source drivers via side-effect
// imports. Open does not accept a context or per-call timeout; callers that need
// bounded validation should use [Check], which intentionally avoids opening
// GitHub sources during health checks.
func Open(sourceURL string) (source.Driver, error) {
	return source.Open(sourceURL)
}
