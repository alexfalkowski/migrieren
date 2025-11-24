package source

import (
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

// Open opens a source driver.
func Open(sourceURL string) (source.Driver, error) {
	return source.Open(sourceURL)
}
