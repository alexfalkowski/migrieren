// Package github provides a golang-migrate source driver that reads migration
// files from a GitHub repository.
//
// It exists so Migrieren depends on a maintained go-github release
// (github.com/google/go-github/v72) instead of the upstream golang-migrate
// GitHub source driver, which pins github.com/google/go-github/v39 and
// transitively imports the unmaintained golang.org/x/crypto/openpgp package
// (advisory GO-2026-5932).
//
// The driver registers itself under the "github" source name in an init
// function, so importing this package for its side effects replaces the
// upstream driver. Do not also import
// github.com/golang-migrate/migrate/v4/source/github, or source.Register will
// panic on the duplicate name.
package github
