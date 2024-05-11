package main

import (
	"os"

	sc "github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/migrieren/cmd"
)

func main() {
	if err := command().Run(); err != nil {
		os.Exit(1)
	}
}

func command() *sc.Command {
	c := sc.New(cmd.Version)
	c.RegisterInput("env:MIGRIEREN_CONFIG_FILE")
	c.AddServer(cmd.ServerOptions...)
	c.AddClient(cmd.ClientOptions...)

	return c
}
