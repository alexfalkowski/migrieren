package migrate

import (
	"fmt"
	"strings"
)

type logger struct {
	logs []string
}

func (l *logger) Printf(format string, v ...any) {
	l.logs = append(l.logs, strings.TrimSpace(fmt.Sprintf(format, v...)))
}

func (l *logger) Verbose() bool {
	return true
}
