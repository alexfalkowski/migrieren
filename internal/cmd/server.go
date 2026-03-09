package cmd

import (
	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// RegisterServer registers the "server" command on command.
//
// The command starts the migrieren service using [Module]. It also registers
// an input flag placeholder (empty default) so runtime config input can be
// provided by callers (for example via `-i`).
func RegisterServer(command cli.Commander) {
	cmd := command.AddServer("server", "Start migrieren server", Module)

	cmd.AddInput(strings.Empty)
}
