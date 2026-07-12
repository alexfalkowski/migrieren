// Package source owns migration source driver wiring and validation.
//
// It registers the supported golang-migrate source drivers and provides a
// bounded health-check path that avoids opening GitHub sources, because the
// registered GitHub driver's Open method is not request-scoped.
package source
