package logger

import (
	"fmt"
	"strings"
	"sync"
)

// New constructs an empty in-memory migration [Logger].
func New() *Logger {
	return &Logger{logs: make([]string, 0)}
}

// Logger records migration logs in memory.
//
// It implements the logger contract used by golang-migrate and is safe for
// concurrent use.
type Logger struct {
	logs []string
	mu   sync.RWMutex
}

// Printf formats and appends a log entry.
//
// Leading and trailing whitespace is removed before storing the message.
func (l *Logger) Printf(format string, v ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logs = append(l.logs, strings.TrimSpace(fmt.Sprintf(format, v...)))
}

// Logs returns the currently recorded log entries.
//
// The returned slice aliases internal storage and should be treated as
// read-only by callers.
func (l *Logger) Logs() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.logs
}

// Verbose reports whether verbose logging is enabled.
//
// This logger is always verbose.
func (l *Logger) Verbose() bool {
	return true
}
