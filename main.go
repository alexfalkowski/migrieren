package main

import (
	"os"

	scmd "github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/migrieren/cmd"
)

func main() {
	command := scmd.New()

	command.AddServer(cmd.ServerOptions)
	command.AddWorker(cmd.WorkerOptions)
	command.AddVersion(cmd.Version)

	if err := command.Run(); err != nil {
		os.Exit(1)
	}
}
