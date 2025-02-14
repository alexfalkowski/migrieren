package main

import (
	sc "github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/migrieren/internal/cmd"
)

func main() {
	command().ExitOnError()
}

func command() *sc.Command {
	command := sc.New(env.NewVersion().String())
	command.RegisterInput(command.Root(), "env:MIGRIEREN_CONFIG_FILE")

	cmd.RegisterServer(command)

	return command
}
