// Package source owns migration source driver wiring and validation.
//
// It registers the supported golang-migrate source drivers and provides a
// bounded health-check path that avoids opening GitHub sources, because the
// upstream GitHub driver does not expose a request-scoped timeout hook.
package source
