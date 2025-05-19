package main

import (
	"context"

	sc "github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/migrieren/internal/cmd"
)

func main() {
	command().ExitOnError(context.Background())
}

func command() *sc.Command {
	fs := os.NewFS()
	command := sc.New(env.NewName(fs), env.NewVersion())

	cmd.RegisterServer(command)

	return command
}
