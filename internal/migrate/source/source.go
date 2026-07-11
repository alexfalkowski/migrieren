package source

import (
	"os"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	_ "github.com/alexfalkowski/migrieren/internal/migrate/source/github"
	"github.com/alexfalkowski/migrieren/internal/migrate/url"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// VersionSource traverses the ordered versions available from a migration
// source. An opened golang-migrate source driver satisfies this interface.
type VersionSource interface {
	First() (uint, error)
	Next(version uint) (uint, error)
}

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
//
// Callers own the returned driver and must close it after successful opens.
func Open(sourceURL string) (source.Driver, error) {
	return source.Open(sourceURL)
}

// Versions returns all migration versions from driver in ascending order.
//
// An empty source returns an empty slice. Errors other than the expected
// os.ErrNotExist end-of-list signal are returned unchanged.
func Versions(driver VersionSource) ([]uint64, error) {
	version, err := driver.First()
	versions := make([]uint64, 0)
	for {
		if errors.Is(err, os.ErrNotExist) {
			return versions, nil
		}
		if err != nil {
			return nil, err
		}

		versions = append(versions, uint64(version))

		version, err = driver.Next(version)
	}
}
