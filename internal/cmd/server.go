package cmd

import (
	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// RegisterServer adds the `server` command to command.
//
// The command boots the Migrieren service using the shared server module.
func RegisterServer(command cli.Commander) {
	cmd := command.AddServer("server", "Start migrieren server", Module)

	cmd.AddInput(strings.Empty)
}
