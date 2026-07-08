package logger

import (
	"fmt"
	"strings"
	"sync"
)

// truncated is the prefix of the marker entry emitted when log lines are dropped.
const truncated = "migration logs truncated"

// New returns an in-memory logger for migration output that returns at most
// limit entries.
//
// limit must be positive; defaulting is the caller's responsibility (see
// [github.com/alexfalkowski/migrieren/internal/migrate.Logs.GetMax]). When the
// captured output exceeds limit, the oldest lines are dropped so the most recent
// lines, which are closest to a failure, are retained.
func New(limit int) *Logger {
	return &Logger{limit: limit}
}

// Logger captures the most recent migration log lines in memory, bounded to a
// caller-supplied limit.
//
// It satisfies the logging interface expected by golang-migrate and is safe for
// concurrent use.
type Logger struct {
	logs    []string
	dropped int
	limit   int
	mu      sync.RWMutex
}

// Printf formats and stores a log line, retaining only the most recent lines
// within the configured limit.
//
// Once any line is dropped, one slot of the limit is reserved for the truncation
// marker added by [Logger.Logs], so the returned slice never exceeds limit
// entries. Dropping reslices the retained tail instead of rebuilding the slice
// on every call, so capture stays cheap for large migrations.
func (l *Logger) Printf(format string, v ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logs = append(l.logs, strings.TrimSpace(fmt.Sprintf(format, v...)))

	budget := l.limit
	if l.dropped > 0 {
		budget = l.limit - 1
	}
	if len(l.logs) <= budget {
		return
	}

	keep := l.limit - 1
	drop := len(l.logs) - keep
	l.logs = l.logs[drop:]
	l.dropped += drop
}

// Logs returns the captured log lines in insertion order.
//
// When earlier lines were dropped, the first entry is a marker that begins with
// "migration logs truncated" and reports how many recent lines are shown out of
// the total captured, for example "migration logs truncated (showing last 99 of
// 512)".
func (l *Logger) Logs() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.dropped == 0 {
		return append([]string(nil), l.logs...)
	}

	marker := fmt.Sprintf("%s (showing last %d of %d)", truncated, len(l.logs), l.dropped+len(l.logs))
	logs := make([]string, 0, len(l.logs)+1)
	logs = append(logs, marker)

	return append(logs, l.logs...)
}

// Verbose reports that verbose migration logging is enabled.
func (l *Logger) Verbose() bool {
	return true
}
