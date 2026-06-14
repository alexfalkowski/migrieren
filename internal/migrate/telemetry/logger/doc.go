// Package logger captures golang-migrate log output for API responses.
//
// The logger is concurrency-safe, bounded in memory, and marks truncation so
// transports can return recent migration logs without exposing unbounded output.
package logger
