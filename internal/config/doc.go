// Package config defines Migrieren's top-level runtime configuration.
//
// It combines the shared go-service runtime config with service-specific health
// and migration sections, then exposes projection helpers used by dependency
// injection after validation has ensured the required sections are present.
package config
