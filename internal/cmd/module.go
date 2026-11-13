package cmd

import (
	"github.com/alexfalkowski/go-service/v3/module"
	"github.com/alexfalkowski/go-signal/v2"
	v1 "github.com/alexfalkowski/migrieren/internal/api/v1"
	"github.com/alexfalkowski/migrieren/internal/config"
	"github.com/alexfalkowski/migrieren/internal/health"
)

// Module wires the full server command dependency graph into the DI container.
var Module = signal.Module(
	module.Server,
	config.Module,
	health.Module,
	v1.Module,
)
