package main

import (
	sc "github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/migrieren/internal/cmd"
)

func main() {
	command().ExitOnError()
}

func command() *sc.Command {
	c := sc.New(cmd.Version)
	c.RegisterInput(c.Root(), "env:MIGRIEREN_CONFIG_FILE")
	c.AddServer("server", "Start migrieren server", cmd.ServerOptions...)

	return c
}
