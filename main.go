package main

import (
	"github.com/alexfalkowski/go-service/v3/cli"
	"github.com/alexfalkowski/go-service/v3/context"
	"github.com/alexfalkowski/go-service/v3/os"
	"github.com/alexfalkowski/migrieren/v2/internal/cmd"
)

var app = cli.NewApplication(func(command cli.Commander) {
	cmd.RegisterServer(command)
})

func main() {
	os.Exit(app.RunCode(context.Background()))
}
