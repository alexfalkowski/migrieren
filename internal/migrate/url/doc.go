// Package url parses and normalizes migration-related URLs.
//
// It rejects empty or schemeless values for configured source/database URLs and
// converts Migrieren's pgx5 database URLs into the PostgreSQL URL shape consumed
// by the underlying SQL driver.
package url
