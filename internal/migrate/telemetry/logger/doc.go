// Package logger provides an in-memory migrate logger implementation.
//
// It satisfies the logging interface expected by
// github.com/golang-migrate/migrate/v4 and is used to capture migration output
// so callers can return logs in API responses.
//
// Logger is safe for concurrent use via internal locking.
package logger
