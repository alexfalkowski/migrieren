package cmd

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/module"
	v1 "github.com/alexfalkowski/migrieren/internal/api/v1"
	"github.com/alexfalkowski/migrieren/internal/config"
	"github.com/alexfalkowski/migrieren/internal/health"
)

// Module for fx.
var Module = di.Module(
	module.Server,
	config.Module,
	health.Module,
	v1.Module,
)
