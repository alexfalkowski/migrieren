package cmd

import (
	"github.com/alexfalkowski/go-service/v2/module"
	v1 "github.com/alexfalkowski/migrieren/internal/api/v1"
	"github.com/alexfalkowski/migrieren/internal/config"
	"github.com/alexfalkowski/migrieren/internal/health"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	module.Server,
	config.Module,
	health.Module,
	v1.Module,
)
