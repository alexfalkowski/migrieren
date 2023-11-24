package main

import (
	"os"

	scmd "github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/migrieren/cmd"
)

func main() {
	if err := command().Run(); err != nil {
		os.Exit(1)
	}
}

func command() *scmd.Command {
	command := scmd.New(cmd.Version)

	command.AddServer(cmd.ServerOptions...)
	command.AddClient(cmd.ClientOptions...)

	return command
}
