package logger

import (
	"fmt"
	"strings"
	"sync"
)

// New returns an in-memory logger for migration output.
func New() *Logger {
	return &Logger{logs: make([]string, 0)}
}

// Logger captures migration log lines in memory.
//
// It satisfies the logging interface expected by golang-migrate and is safe for
// concurrent use.
type Logger struct {
	logs []string
	mu   sync.RWMutex
}

// Printf formats and stores a log line.
func (l *Logger) Printf(format string, v ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logs = append(l.logs, strings.TrimSpace(fmt.Sprintf(format, v...)))
}

// Logs returns the captured log lines in insertion order.
func (l *Logger) Logs() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.logs
}

// Verbose reports that verbose migration logging is enabled.
func (l *Logger) Verbose() bool {
	return true
}
