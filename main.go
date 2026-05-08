package main

import (
	"github.com/alexfalkowski/go-service/v2/cli"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/migrieren/internal/cmd"
)

var app = cli.NewApplication(func(command cli.Commander) {
	cmd.RegisterServer(command)
})

func main() {
	os.Exit(app.RunCode(context.Background()))
}
