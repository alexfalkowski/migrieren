package source

import (
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

// Open opens a migration source driver for sourceURL.
//
// Supported schemes depend on registered source drivers (see package docs).
// Underlying driver errors are returned unchanged.
func Open(sourceURL string) (source.Driver, error) {
	return source.Open(sourceURL)
}
