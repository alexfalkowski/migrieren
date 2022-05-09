package migrate

import (
	"fmt"
)

type logger struct {
	logs []string
}

func (l *logger) Printf(format string, v ...any) {
	l.logs = append(l.logs, fmt.Sprintf(format, v...))
}

func (l *logger) Verbose() bool {
	return true
}
