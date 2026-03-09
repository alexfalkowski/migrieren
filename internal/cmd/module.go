package cmd

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/module"
	v1 "github.com/alexfalkowski/migrieren/internal/api/v1"
	"github.com/alexfalkowski/migrieren/internal/config"
	"github.com/alexfalkowski/migrieren/internal/health"
)

// Module registers the dependencies required by the "server" CLI command.
//
// It composes runtime server infrastructure, service configuration, health
// checks, and API transport wiring.
var Module = di.Module(
	module.Server,
	config.Module,
	health.Module,
	v1.Module,
)
