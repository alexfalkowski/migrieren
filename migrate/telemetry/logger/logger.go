package logger

import (
	"fmt"
	"strings"
	"sync"
)

// New used for migrations.
func New() *Logger {
	return &Logger{logs: make([]string, 0)}
}

// Logger used for migrations.
type Logger struct {
	logs []string
	mu   sync.RWMutex
}

// Printf the log message.
func (l *Logger) Printf(format string, v ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logs = append(l.logs, strings.TrimSpace(fmt.Sprintf(format, v...)))
}

// Logs that have been written.
func (l *Logger) Logs() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.logs
}

// Verbose for our logger.
func (l *Logger) Verbose() bool {
	return true
}
